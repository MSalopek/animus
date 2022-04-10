package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TODO: abstract the DB logic into a set of interfaces

const defaultSecret = "pleaseDontUsethisstring"

type HttpAPI struct {
	engine *gin.Engine
	server *http.Server
	db     *gorm.DB
	secret string
	logger *logrus.Logger

	done chan struct{}
	port string
}

func New(port string, db *gorm.DB, log *logrus.Logger, done chan struct{}) *HttpAPI {
	if db == nil {
		panic("DB must be provided")
	}

	engine := gin.New()
	s := &HttpAPI{
		engine: engine,
		server: &http.Server{
			Addr:    port,
			Handler: engine,
		},
		db:     db,
		secret: defaultSecret,
		done:   done,
		logger: log,
	}

	engine.Use(requestLogger(log))
	return s
}

func (api *HttpAPI) registerHandlers() {
	root := api.engine.Group("/api")

	root.GET("/ping", WIPresponder)
	root.GET("/login/", api.Login)
	root.POST("/register", api.Register)

	root.POST("/files/add", WIPresponder)
	root.POST("/files/pin/:id", WIPresponder)
	root.DELETE("/files/delete/:id", WIPresponder)
	root.PUT("/files/update/:id", WIPresponder)

	// get single file/directory metadata
	root.GET("/files/:id/stat")

	// get paginated list of user's files/directories
	root.GET("/files/user/:id", WIPresponder)

	root.POST("/gates/", WIPresponder)
	root.GET("/gates/:id", WIPresponder)
	root.DELETE("/gates/:id", WIPresponder)
	root.PATCH("/gates/:id", WIPresponder)

	// get all gates for user
	root.GET("/gates/user/:id", WIPresponder)
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
