package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/engine/api/auth"
	"github.com/msalopek/animus/engine/repo"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func handleCORS(allowedOrigins []string) gin.HandlerFunc {
	if len(allowedOrigins) == 0 {
		panic("missing allowed origins")
	}
	return cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{
			"GET",
			"PUT",
			"POST",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Authorization",
			"Cache-Control",
			"Content-Type",
			"Accept-Encoding",
			"X-CSRF-Token",
			"X-Requested-With",
			"Origin",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
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

func authorizeUser(repo *repo.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// injected by authorizeRequest
		email := c.GetString("email")
		if email == "" {
			abortWithError(c, http.StatusUnauthorized, engine.ErrUnauthorized)
			return
		}

		user, err := repo.GetUserByEmail(email)
		if err != nil {
			abortWithError(c, http.StatusForbidden, engine.ErrForbidden)
			return
		}

		// inject userID into gin.Context
		c.Set("userID", int(user.ID))
		c.Next()
	}
}
