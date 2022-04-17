package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/msalopek/animus/engine/api/auth"
	"github.com/msalopek/animus/engine/repo"
	log "github.com/sirupsen/logrus"
)

// TODO: abstract the DB logic into a set of interfaces

const defaultSecret = "pleaseDontUsethisstring"
const defaultExpiration = 1 * time.Hour

type HttpAPI struct {
	engine *gin.Engine
	server *http.Server
	ipfs   *shell.Shell
	repo   *repo.Repo
	auth   *auth.Auth
	logger *log.Logger

	done chan struct{}
	port string
}

func New(port string, ipfsURL string, repo *repo.Repo, logger *log.Logger, done chan struct{}) *HttpAPI {
	if repo == nil {
		panic("Repo must be provided")
	}

	engine := gin.New()
	shell := shell.NewShell(ipfsURL)
	s := &HttpAPI{
		engine: engine,
		server: &http.Server{
			Addr:    port,
			Handler: engine,
		},
		repo: repo,
		ipfs: shell,
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
	auth.POST("/storage/add", api.UploadFile)
	auth.GET("/user/:id/storage", api.GetUserUploads)
	// auth.GET("/storage/:cid", WIPresponder)

	// TODO
	// auth.GET("/storage/ls/:cid", api.ProxyCommandLs)
	// auth.POST("/storage/pin/:id", WIPresponder)
	// auth.DELETE("/storage/delete/:id", WIPresponder)

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

	api.logger.Info("graceful shutdown successful")
	return nil
}

func (api *HttpAPI) IPFSUpload(r io.Reader) (string, error) {
	fileHash, err := api.ipfs.Add(r)
	if err != nil {
		return "", err
	}
	return fileHash, nil
}
