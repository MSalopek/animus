package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine/api/auth"
	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/storage"
	log "github.com/sirupsen/logrus"
)

const defaultSecret = "pleaseDontUsethisstring"
const defaultExpiration = 1 * time.Hour

type Publisher interface {
	Publish(topic string, body []byte) error
	Stop()
}

type Config struct {
	AuthSecret          string `json:"auth_secret,omitempty" yaml:"auth_secret,omitempty"`
	AuthAuthority       string `json:"auth_authority,omitempty" yaml:"auth_authority,omitempty"`
	AuthExpirationHours string `json:"auth_expiration_hours,omitempty" yaml:"auth_expiration_hours,omitempty"`
	HttpPort            string `json:"http_port,omitempty" yaml:"http_port,omitempty"`
	LogLevel            string `json:"log_level,omitempty" yaml:"log_level,omitempty"`

	Bucket  string         `json:"bucket,omitempty" yaml:"bucket,omitempty"`
	Storage storage.Config `json:"storage,omitempty" yaml:"storage,omitempty"`
}

type HttpAPI struct {
	engine    *gin.Engine
	server    *http.Server
	storage   *storage.Manager
	repo      *repo.Repo
	auth      *auth.Auth
	logger    *log.Logger
	publisher Publisher

	done chan struct{}

	cfg *Config
}

func New(cfg *Config, repo *repo.Repo, logger *log.Logger, done chan struct{}) *HttpAPI {
	if cfg == nil {
		panic("config must be provided")
	}
	if repo == nil {
		panic("Repo must be provided")
	}

	engine := gin.New()
	s := &HttpAPI{
		engine: engine,
		server: &http.Server{
			Addr:    cfg.HttpPort,
			Handler: engine,
		},
		repo: repo,
		auth: &auth.Auth{
			Secret:          defaultSecret,
			Authority:       "AnimusEngine",
			ExpirationHours: defaultExpiration,
		},
		done:   done,
		logger: logger,
	}

	engine.Use(requestLogger(logger))
	return s
}

func (api *HttpAPI) registerHandlers() {
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
	auth.POST("/storage/add-dir", api.UploadFile)
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

func (api *HttpAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.engine.ServeHTTP(w, r)
}

func (api *HttpAPI) Start() error {
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

func (api *HttpAPI) Stop() error {
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

	api.publisher.Stop()
	api.logger.Info("graceful shutdown successful")
	return nil
}
