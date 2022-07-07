package api

import (
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/model"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// APIClientProvider fetches data about API Client that matches provided key.
type APIClientProvider interface {
	GetApiClientByKey(key string) (*model.APIClient, error)
}

func requestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.WithFields(
			log.Fields{
				"method":  c.Request.Method,
				"URL":     c.Request.URL,
				"headers": c.Request.Header,
			}).Info("new api request")
		c.Next()
	}
}

// client is an application using client_key + client_secret
func authorizeClientRequest(prov APIClientProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		if bearer != "" {
			abortWithError(c, http.StatusBadRequest, engine.ErrInvalidClientAuth)
			return
		}

		key := c.Request.Header.Get("X-API-KEY")
		if key == "" {
			abortWithError(c, http.StatusUnauthorized, engine.ErrInvalidClientAuth)
			return
		}

		sig := c.Request.Header.Get("X-API-SIGN")
		if sig == "" {
			abortWithError(c, http.StatusUnauthorized, engine.ErrInvalidClientAuth)
			return
		}

		decSig, err := hex.DecodeString(sig)
		if err != nil {
			abortWithError(c, http.StatusBadRequest, engine.ErrInvalidClientSignature)
			return
		}

		client, err := prov.GetApiClientByKey(key)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			abortWithError(c, http.StatusUnauthorized, engine.ErrInvalidClientAuth)
			return
		} else if err != nil {
			abortWithError(c, http.StatusInternalServerError, engine.ErrInternalError)
			return
		}

		if !ValidateSignature([]byte(c.Request.RequestURI),
			decSig,
			[]byte(client.ClientSecret)) {

			abortWithError(c, http.StatusUnauthorized, engine.ErrInvalidClientAuth)
			return
		}

		// make available to handlers though context
		c.Set("userID", int(client.UserID))
		c.Set("email", client.Email)
		c.Next()
	}
}
