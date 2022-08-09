package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/pinner"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

var cfgPath = flag.String("config", "", "path to pinnerd configuration file")

func main() {
	flag.Parse()
	if *cfgPath == "" {
		log.Fatal("config path not provided")
	}
	file, err := ioutil.ReadFile(*cfgPath)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	var cfg pinner.Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatalf("error unmarshaling config %s", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	ctx := exitSignal()

	logger := log.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	if cfg.TextLogs {
		logger.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	}

	if cfg.Debug {
		logger.SetLevel(log.DebugLevel)
		// buf, _ := yaml.Marshal(cfg)
		// logger.Debug(fmt.Sprintf("-- config --\n%s-- ------ --", string(buf)))
	}

	var db *gorm.DB
	if cfg.DbDSN != "" {
		db, err = gorm.Open(postgres.Open(cfg.DbDSN), &gorm.Config{})
		if err != nil {
			log.Fatal("could not connect to database - terminating")
		}
	} else {
		db, err = gorm.Open(postgres.Open(cfg.DbURI), &gorm.Config{})
		if err != nil {
			log.Fatal("could not connect to database - terminating")
		}
	}
	repo := repo.New(db)

	pinner := pinner.New(&cfg, repo, logger)
	wg.Add(1)
	// go queue.WgChanBodyLogger(&wg, pinner.Messages)
	go pinner.HandlePinRequests(&wg)

	wg.Add(1)
	go pinner.HandleWebhooks(&wg)

	<-ctx.Done()
	logger.Info("pinner stopping")
	pinner.Stop()
	wg.Wait()
	logger.Info("pinner stopped")
}

// exitSignal cancels the context on SIGINT, SIGTERM
// use the returned context as the root context
func exitSignal() context.Context {
	ctx, stop := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		stop()
	}()
	return ctx
}
