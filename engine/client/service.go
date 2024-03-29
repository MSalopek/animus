package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/engine/repo"
	"github.com/msalopek/animus/queue"
	"github.com/msalopek/animus/storage"
	log "github.com/sirupsen/logrus"
)

// same as 100*1024*1024 -> 100MB
const defaultMaxBodySize = 100 << 20

type ClientAPI struct {
	sh      *shell.Shell
	engine  *gin.Engine
	server  *http.Server
	storage *storage.Manager
	repo    *repo.Repo
	logger  *log.Logger

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

	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()
	s := &ClientAPI{
		cfg: cfg,

		engine: e,
		server: &http.Server{
			Addr:    cfg.HttpPort,
			Handler: e,
		},

		sh:        shell.NewShell(cfg.NodeApiURL),
		repo:      repo,
		storage:   storage.MustNewManager(cfg.Storage),
		publisher: queue.MustNewPublisher(cfg.NsqPinnerTopic, cfg.NsqdURL),
		done:      done,
		logger:    logger,
	}

	e.Use(engine.LogRequest(logger))
	e.MaxMultipartMemory = defaultMaxBodySize
	return s
}

func (api *ClientAPI) registerHandlers() {
	root := api.engine.Group("/")
	root.GET("/heartbeat", api.Heartbeat)

	auth := root.Group("/auth").Use(
		authorizeClientRequest(api.repo),
		checkBodySize(),
	)
	auth.GET("/whoami", api.WhoAmI)

	auth.GET("/storage", api.GetStorageRecords)
	auth.GET("/storage/id/:id", api.GetStorageRecord)
	auth.DELETE("/storage/id/:id", api.DeleteStorageRecord)
	auth.GET("/storage/cid/:cid", api.GetStorageRecordByCid)

	auth.POST("/storage/sync-add-file", api.SyncUploadFile)
	auth.POST("/storage/add-file", api.UploadFile)
	auth.POST("/storage/add-dir", api.UploadDir)
	auth.POST("/storage/pin/id/:id", api.RequestPin)
	auth.POST("/storage/unpin/id/:id", api.RequestUnpin)

	auth.PATCH("/user", api.UpdateUser)
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

func (api *ClientAPI) Heartbeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"service":   "client_api",
		"timestamp": time.Now().UnixNano(),
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
