package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/engine/api/auth"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

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

func authorizeRequest(secret, autority string) gin.HandlerFunc {
	auth := auth.Auth{
		Secret:    secret,
		Authority: autority,
	}

	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		if bearer == "" {
			abortWithError(c, 403, engine.ErrNoAuthHeader)
			return
		}

		ts := strings.Split(bearer, "Bearer ")

		if len(ts) == 2 {
			bearer = strings.TrimSpace(ts[1])
		} else {
			abortWithError(c, 400, engine.ErrInvalidAuthToken)
			return
		}

		claims, err := auth.ValidateToken(bearer)
		if err != nil {
			// TODO: log exact error
			abortWithError(c, 401, engine.ErrUnauthorized)
			return
		}

		// inject email into gin.Context
		c.Set("email", claims.Email)
		c.Next()

	}
}
