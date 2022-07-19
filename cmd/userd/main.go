package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/engine/user"
)

var cfgPath = flag.String("config", "", "path to userd configuration file")

func main() {
	flag.Parse()
	if *cfgPath == "" {
		log.Fatal("config path not provided")
	}
	file, err := ioutil.ReadFile(*cfgPath)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	var cfg user.Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatalf("error unmarshaling config %s", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	ctx := exitSignal()
	done := make(chan struct{})
	logger := log.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	if cfg.TextLogs {
		logger.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
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

	svc := user.New(&cfg, repo, logger, done)

	go svc.Start()

	<-ctx.Done()
	err = svc.Stop()
	if err != nil {
		logger.Error(err)
	}
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
