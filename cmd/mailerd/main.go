package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/msalopek/animus/mailer"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var cfgPath = flag.String("config", "", "path to clientd configuration file")

func main() {
	flag.Parse()
	if *cfgPath == "" {
		log.Fatal("config path not provided")
	}
	file, err := ioutil.ReadFile(*cfgPath)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	var cfg mailer.Config
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

	if *&cfg.Debug {
		logger.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
		logger.SetLevel(log.DebugLevel)
	}

	m := mailer.New(
		cfg.ApiKey,
		cfg.ApiSecret,
		cfg.Topic,
		cfg.NsqLookup,
		cfg.Sender,
		cfg.SenderName,
		logger,
		cfg.DryRun,
	)
	// "matija@animus.store", "Matija from Animus",
	wg.Add(1)
	go m.HandleMessages(&wg)

	<-ctx.Done()
	logger.Info("mailer stopping")
	m.Stop()
	wg.Wait()
	logger.Info("mailer stopped")
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
