package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/msalopek/animus/engine/api"
	"github.com/msalopek/animus/engine/repo"
)

const defaultHttpPort = ":8083"
const defaultIPFSApi = "localhost:5001"
const defaultStoragePath = "./storage"
const defaultDBUri = "postgres://animus:animus@localhost:5432/animus"
const defaultDSN = "host=localhost user=animus password=animus dbname=animus port=5432 sslmode=disable"

func main() {
	ctx := exitSignal()

	db, err := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{})
	if err != nil {
		panic("could not connect to database - terminating")
	}
	repo := repo.New(db)

	done := make(chan struct{})
	logger := log.New()
	logger.Out = os.Stdout
	logger.SetFormatter(&log.TextFormatter{})

	api := api.New(
		defaultHttpPort,
		defaultIPFSApi,
		repo,
		logger,
		done,
	)

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
