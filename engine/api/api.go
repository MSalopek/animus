package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TODO: abstract the DB logic into a set of interfaces

const defaultSecret = "pleaseDontUsethisstring"

type HttpAPI struct {
	engine *gin.Engine
	server *http.Server
	ipfs   *shell.Shell
	db     *gorm.DB
	secret string
	logger *logrus.Logger

	done chan struct{}
	port string
}

func New(port string, ipfsURL string, db *gorm.DB, log *logrus.Logger, done chan struct{}) *HttpAPI {
	if db == nil {
		panic("DB must be provided")
	}

	engine := gin.New()
	shell := shell.NewShell(ipfsURL)
	s := &HttpAPI{
		engine: engine,
		server: &http.Server{
			Addr:    port,
			Handler: engine,
		},
		db:     db,
		ipfs:   shell,
		secret: defaultSecret,
		done:   done,
		logger: log,
	}

	engine.Use(requestLogger(log))
	return s
}

func (api *HttpAPI) registerHandlers() {
	root := api.engine.Group("/api")

	root.GET("/ping", api.Ping)
	root.GET("/login/", api.Login)
	root.POST("/register", api.Register)

	auth := root.Group("/auth").Use(authorizeRequest(
		api.secret,
		"AnimusEngine",
	))
	auth.GET("/whoami", api.WhoAmI)
	// get paginated list of user's files/directories
	auth.GET("/manager/user/:id", WIPresponder)
	auth.POST("/manager/add", WIPresponder)
	// auth.POST("/manager/pin/:id", WIPresponder)
	// auth.DELETE("/manager/delete/:id", WIPresponder)
	// auth.PUT("/manager/update/:id", WIPresponder)

	// // get single file/directory metadata
	// auth.GET("/manager/:id/stat")

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
