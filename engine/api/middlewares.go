package api

import (
	"encoding/base64"
	"errors"
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
	"gorm.io/gorm"
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
			"X-API-KEY",
			"X-API-SIG",
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

// user is an person using email + password
// the user must first login and acquire a JWT token to be able to proceed
func authorizeUserRequest(authority *auth.Auth) gin.HandlerFunc {
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

// client is an application using client_key + client_secret
func authorizeClientRequest(repo *repo.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get("X-API-KEY")
		if key == "" {
			abortWithError(c, http.StatusForbidden, engine.ErrInvalidClientAuth)
			return
		}

		sig := c.Request.Header.Get("X-API-SIGN")
		if sig == "" {
			abortWithError(c, http.StatusForbidden, engine.ErrInvalidClientAuth)
			return
		}

		client, err := repo.GetApiClientByKey(key)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			abortWithError(c, http.StatusForbidden, engine.ErrForbidden)
			return
		} else if err != nil {
			abortWithError(c, http.StatusInternalServerError, engine.ErrInternalError)
			return
		}

		decSig, err := base64.StdEncoding.DecodeString(sig)
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, engine.ErrInternalError)
			return
		}
		valid := ValidateSignature([]byte(c.Request.URL.String()), []byte(decSig), []byte(client.ClientSecret))
		if !valid {
			abortWithError(c, http.StatusUnauthorized, err.Error())
			return
		}

		// inject email into gin.Context
		c.Set("userID", client.UserID)
		c.Set("email", client.Email)
		c.Next()
	}
}
