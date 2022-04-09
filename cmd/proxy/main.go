package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"

	"github.com/msalopek/animus/api"
)

const defaultHttpPort = ":8083"
const defaultIPFSApi = "localhost:50001"
const defaultStoragePath = "./storage"
const defaultDBUri = "postgres://animus:animus@localhost:5436/animus"

func main() {
	ctx := exitSignal()

	dbPool, err := pgxpool.Connect(context.Background(), defaultDBUri)
	if err != nil {
		panic("could not connect to database - terminating")
	}
	defer dbPool.Close()

	done := make(chan struct{})
	logger := log.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&log.TextFormatter{})

	api := api.New(defaultHttpPort, dbPool, logger, done)

	go api.Start()

	<-ctx.Done()
	err = api.Stop()
	if err != nil {
		fmt.Println(err)
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
