package api

import (
	"net/http"
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

func authorizeRequest(authority *auth.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		if bearer == "" {
			abortWithError(c, http.StatusForbidden, engine.ErrNoAuthHeader)
			return
		}

		ts := strings.Split(bearer, "Bearer ")

		if len(ts) == 2 {
			bearer = strings.TrimSpace(ts[1])
		} else {
			abortWithError(c, http.StatusBadRequest, engine.ErrInvalidAuthToken)
			return
		}

		claims, err := authority.ValidateToken(bearer)
		if err != nil {
			// TODO: log exact error
			// abortWithError(c, http.StatusUnauthorized, engine.ErrUnauthorized)
			abortWithError(c, http.StatusUnauthorized, err.Error())
			return
		}

		// inject email into gin.Context
		c.Set("email", claims.Email)
		c.Next()
	}
}
