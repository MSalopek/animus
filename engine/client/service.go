package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/queue"
	"github.com/msalopek/animus/storage"
	log "github.com/sirupsen/logrus"
)

type ClientAPI struct {
	engine    *gin.Engine
	server    *http.Server
	storage   *storage.Manager
	repo      *repo.Repo
	logger    *log.Logger
	publisher queue.Publisher

	done chan struct{}

	cfg *Config
}

func New(cfg *Config, repo *repo.Repo, logger *log.Logger, done chan struct{}) *ClientAPI {
	if cfg == nil {
		panic("config must be provided")
	}
	if repo == nil {
		panic("Repo must be provided")
	}

	e := gin.New()
	s := &ClientAPI{
		cfg: cfg,

		engine: e,
		server: &http.Server{
			Addr:    cfg.HttpPort,
			Handler: e,
		},
		repo:      repo,
		storage:   storage.MustNewManager(cfg.Storage),
		publisher: queue.MustNewPublisher(cfg.NsqTopic, cfg.NsqdURL),
		done:      done,
		logger:    logger,
	}

	e.Use(engine.LogRequest(logger))
	return s
}

func (api *ClientAPI) registerHandlers() {
	root := api.engine.Group("/")
	root.GET("/ping", api.Ping)

	auth := root.Group("/auth").Use(
		authorizeClientRequest(api.repo),
	)
	auth.GET("/whoami", api.WhoAmI)

	auth.GET("/storage", api.GetStorageRecords)
	auth.GET("/storage/id/:id", api.GetStorageRecord)
	auth.POST("/storage/add-file", api.UploadFile)
	auth.POST("/storage/add-dir", api.UploadDir)

	// TODO:
	// auth.POST("/storage/pin/:id", WIPresponder)
	// auth.POST("/storage/unpin/:id", WIPresponder)
	// auth.GET("/storage/pin/status/:id/", WIPresponder)
}

func (api *ClientAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.engine.ServeHTTP(w, r)
}

func (api *ClientAPI) Start() error {
	api.registerHandlers()
	api.logger.Info("starting clientd service")
	err := api.server.ListenAndServe()
	if err == http.ErrServerClosed {
		// wait for the graceful shutdown
		<-api.done
		return nil
	}

	api.logger.WithFields(log.Fields{
		"error": err.Error(),
	}).Error("unexpected error")
	return err
}

func (api *ClientAPI) Stop() error {
	api.logger.Info("graceful shutdown initiated")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer close(api.done)

	api.publisher.Stop()

	// graceful shutdown (serve all in flight)
	// Error from closing listeners, or context timeout:
	err := api.server.Shutdown(ctx)
	if err != nil {
		api.logger.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error shutting down")
		return fmt.Errorf("HTTP server Shutdown: %v", err)
	}

	// api.publisher.Stop()
	api.logger.Info("graceful shutdown successful")
	return nil
}

func (api *ClientAPI) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

// WhoAmI returns data that identifies the API Key user.
// It is essentially the Ping method but for authenticated users.
func (api *ClientAPI) WhoAmI(c *gin.Context) {
	email := c.GetString("email")
	if len(email) < 1 {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	user, err := api.repo.GetUserByEmail(email)
	if err == gorm.ErrRecordNotFound {
		engine.AbortErr(c, http.StatusNotFound, engine.ErrNotFound)
		return
	}

	if err != nil {
		engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
		return
	}

	c.JSON(200, user)
}
