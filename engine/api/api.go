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
	"github.com/msalopek/animus/queue"
	"github.com/msalopek/animus/storage"
	log "github.com/sirupsen/logrus"
)

type AnimusAPI struct {
	engine    *gin.Engine
	server    *http.Server
	storage   *storage.Manager
	repo      *repo.Repo
	auth      *auth.Auth
	logger    *log.Logger
	publisher queue.Publisher

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
	// CORS defined on all routes
	engine.Use(handleCORS(cfg.CORSOrigins))
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
		storage:   storage.MustNewManager(cfg.Storage),
		publisher: queue.MustNewPublisher(cfg.NsqTopic, cfg.NsqdURL),
		done:      done,
		logger:    logger,
	}

	engine.Use(requestLogger(logger))
	return s
}

func (api *AnimusAPI) registerHandlers() {
	root := api.engine.Group("/api")
	root.GET("/ping", api.Ping)
	root.POST("/login", api.Login)
	root.POST("/register", api.Register)

	auth := root.Group("/auth").Use(
		authorizeUserRequest(api.auth),
		authorizeUser(api.repo),
	)
	auth.GET("/whoami", api.WhoAmI)

	auth.GET("/user/keys", api.GetUserKeys)
	auth.POST("/user/keys", api.CreateUserKey)
	auth.PATCH("/user/keys/id/:id", api.GetUserUploads)
	auth.DELETE("/user/keys/id/:id", api.GetUserUploads)

	auth.GET("/user/storage", api.GetUserUploads)
	auth.POST("/storage/add-file", api.UploadFile)
	auth.POST("/storage/add-dir", api.UploadDir)

	// /gate group is for key-secret auth for API access
	gate := root.Group("/gate").Use(
		authorizeClientRequest(api.repo),
	)
	gate.GET("/whoami", api.WhoAmI)
	gate.POST("/storage/add-file", api.UploadFile)
	gate.POST("/storage/add-dir", api.UploadDir)
	gate.GET("/storage/user", api.GetUserUploads)

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
	api.logger.Info("starting animusd service")
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
