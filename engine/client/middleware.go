package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/msalopek/animus/engine"
	"github.com/msalopek/animus/model"

	"gorm.io/gorm"
)

// APIClientProvider fetches data about API Client that matches provided key.
type APIClientProvider interface {
	GetApiClientByKey(key string) (*model.APIClient, error)
}

// client is an application using client_key + client_secret
func authorizeClientRequest(prov APIClientProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		// bearer := c.Request.Header.Get("Authorization")
		// if bearer != "" {
		// 	engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidClientAuth)
		// 	return
		// }
		key := c.Request.Header.Get("X-API-KEY")
		if key == "" {
			engine.AbortErr(c, http.StatusUnauthorized, engine.ErrInvalidClientAuth)
			return
		}

		sig := c.Request.Header.Get("X-API-SIGN")
		if sig == "" {
			engine.AbortErr(c, http.StatusUnauthorized, engine.ErrInvalidClientAuth)
			return
		}

		decSig, err := hex.DecodeString(sig)
		if err != nil {
			engine.AbortErr(c, http.StatusBadRequest, engine.ErrInvalidClientSignature)
			return
		}

		client, err := prov.GetApiClientByKey(key)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			engine.AbortErr(c, http.StatusUnauthorized, engine.ErrInvalidClientAuth)
			return
		} else if err != nil {
			engine.AbortErr(c, http.StatusInternalServerError, engine.ErrInternalError)
			return
		}

		if !ValidateSignature([]byte(c.Request.RequestURI),
			decSig,
			[]byte(client.ClientSecret)) {

			engine.AbortErr(c, http.StatusUnauthorized, engine.ErrInvalidClientAuth)
			return
		}

		// make available to handlers though context
		c.Set("userID", int(client.UserID))
		c.Set("email", client.Email)
		c.Next()
	}
}

// ValidateSignature reports whether signature is a valid HMAC for message.
// https://pkg.go.dev/crypto/hmac
func ValidateSignature(message, signature, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expected := mac.Sum(nil)
	return hmac.Equal(signature, expected)
}

func checkBodySize() gin.HandlerFunc {
	return func(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		c.Request.Body = http.MaxBytesReader(w, c.Request.Body, defaultMaxBodySize)

		c.Next()
	}
}
