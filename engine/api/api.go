package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/engine/api/auth"
	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/storage"
	log "github.com/sirupsen/logrus"
)

type Publisher interface {
	Publish(topic string, body []byte) error
	Stop()
}
type AnimusAPI struct {
	engine    *gin.Engine
	server    *http.Server
	storage   *storage.Manager
	repo      *repo.Repo
	auth      *auth.Auth
	logger    *log.Logger
	publisher Publisher

	done chan struct{}

	cfg *engine.Config
}

func New(cfg *engine.Config, repo *repo.Repo, logger *log.Logger, done chan struct{}) *AnimusAPI {
	if cfg == nil {
		panic("config must be provided")
	}
	if repo == nil {
		panic("Repo must be provided")
	}

	engine := gin.New()
	s := &AnimusAPI{
		cfg: cfg,

		engine: engine,
		server: &http.Server{
			Addr:    cfg.HttpPort,
			Handler: engine,
		},
		repo: repo,
		auth: &auth.Auth{
			Secret:          cfg.AuthSecret,
			Authority:       cfg.AuthAuthority,
			ExpirationHours: time.Duration(cfg.AuthExpirationHours) * time.Hour,
		},
		storage: storage.MustNewManager(cfg.Storage),
		done:    done,
		logger:  logger,
	}

	engine.Use(requestLogger(logger))
	return s
}

func (api *AnimusAPI) registerHandlers() {
	root := api.engine.Group("/api")

	root.GET("/ping", api.Ping)
	root.POST("/login/", api.Login)
	root.POST("/register", api.Register)

	auth := root.Group("/auth").Use(
		authorizeRequest(api.auth),
		authorizeUser(api.repo),
	)
	auth.GET("/whoami", api.WhoAmI)
	auth.POST("/storage/add-file", api.UploadFile)
	auth.POST("/storage/add-dir", api.UploadDir)
	auth.GET("/storage/user", api.GetUserUploads)

	// TODO:
	// auth.POST("/storage/pin/:id", WIPresponder)
	// auth.POST("/storage/unpin/:id", WIPresponder)

	// auth.POST("/gates/", WIPresponder)
	// auth.GET("/gates/:id", WIPresponder)
	// auth.DELETE("/gates/:id", WIPresponder)
	// auth.PATCH("/gates/:id", WIPresponder)

	// // get all gates for user
	// auth.GET("/gates/user/:id", WIPresponder)
}

func (api *AnimusAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.engine.ServeHTTP(w, r)
}

func (api *AnimusAPI) Start() error {
	api.registerHandlers()
	api.logger.Info("starting wallet service")
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

func (api *AnimusAPI) Stop() error {
	api.logger.Info("graceful shutdown initiated")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer close(api.done)

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
