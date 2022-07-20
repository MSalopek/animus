package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/engine/user/auth"
	"github.com/msalopek/animus/queue"
	"github.com/msalopek/animus/storage"
	log "github.com/sirupsen/logrus"
)

type UserAPI struct {
	engine    *gin.Engine
	server    *http.Server
	storage   *storage.Manager
	repo      *repo.Repo
	auth      *auth.Auth
	logger    *log.Logger
	publisher queue.Publisher

	done chan struct{}

	cfg *Config
}

func New(cfg *Config, repo *repo.Repo, logger *log.Logger, done chan struct{}) *UserAPI {
	if cfg == nil {
		panic("config must be provided")
	}
	if repo == nil {
		panic("Repo must be provided")
	}

	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.New()
	e.Use(
		// CORS defined on all routes
		handleCORS(cfg.CORSOrigins),
		engine.LogRequest(logger))
	s := &UserAPI{
		cfg: cfg,

		engine: e,
		server: &http.Server{
			Addr:    cfg.HttpPort,
			Handler: e,
		},
		repo: repo,
		auth: &auth.Auth{
			Secret:          cfg.AuthSecret,
			Authority:       cfg.AuthAuthority,
			ExpirationHours: time.Duration(cfg.AuthExpirationHours) * time.Hour,
		},
		storage:   storage.MustNewManager(cfg.Storage),
		publisher: queue.MustNewPublisher(cfg.NsqPinnerTopic, cfg.NsqdURL),
		done:      done,
		logger:    logger,
	}

	return s
}

func (api *UserAPI) registerHandlers() {
	root := api.engine.Group("/")
	root.GET("/heartbeat", api.Heartbeat)
	root.POST("/login", api.Login)
	root.POST("/register", api.Register)
	root.POST("/activate/email/:email", api.ActivateUser)

	auth := root.Group("/auth").Use(
		authorizeUserRequest(api.auth),
		authorizeUser(api.repo),
	)
	auth.GET("/whoami", api.WhoAmI)

	auth.GET("/user/keys", api.GetUserKeys)
	auth.POST("/user/keys", api.CreateUserKey)
	auth.PATCH("/user/keys/id/:id", api.UpdateUserKey)
	auth.DELETE("/user/keys/id/:id", api.DeleteUserKey)

	auth.GET("/user/storage", api.GetUserUploads)
	auth.GET("/user/storage/id/:id", api.GetStorageRecord)
	auth.DELETE("/user/storage/id/:id", api.DeleteStorageRecord)
	auth.POST("/user/storage/add-file", api.UploadFile)
	auth.POST("/user/storage/add-dir", api.UploadDir)
	auth.GET("/user/storage/download/id/:id", api.UploadDir)

	// IPFS operations
	auth.POST("/user/storage/pin/id/:id", api.RequestPin)
	auth.POST("/user/storage/unpin/id/:id", api.RequestUnpin)

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

func (api *UserAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.engine.ServeHTTP(w, r)
}

func (api *UserAPI) Start() error {
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

func (api *UserAPI) Stop() error {
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

func (api *UserAPI) Heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"service":   "user_api",
		"timestamp": time.Now().UnixNano(),
	})
}
