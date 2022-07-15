package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/msalopek/animus/mailer"
	log "github.com/sirupsen/logrus"
)

var dryRun = flag.Bool("dry-run", false, "output emails to stdout (does not send API requests)")
var key = flag.String("api-key", "", "mailer REST API key")
var secret = flag.String("api-secret", "", "mailer REST API secret")
var topic = flag.String("topic", "", "nsq email topic")
var nsqLookup = flag.String("nsq-lookup", "", "nsq lookup for nsq connections")
var debug = flag.Bool("debug", false, "show debug logs")

func main() {
	flag.Parse()

	if *key == "" {
		log.Fatal("api-key is required")
	}
	if *secret == "" {
		log.Fatal("api-secret is required")
	}
	if *topic == "" {
		log.Fatal("topic is required")
	}
	if *nsqLookup == "" {
		log.Fatal("nsq-lookup is required")
	}

	var wg sync.WaitGroup
	ctx := exitSignal()

	logger := log.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})

	if *debug {
		logger.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
		logger.SetLevel(log.DebugLevel)
	}

	m := mailer.New(*key, *secret, *topic, *nsqLookup, "matija@animus.store", "Matija from Animus", logger, *dryRun)
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
