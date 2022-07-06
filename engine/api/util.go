package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"net/http"

	"github.com/gin-gonic/gin"
)

func abortWithError(c *gin.Context, httpCode int, err string) {
	c.Abort()
	c.JSON(httpCode, gin.H{
		"error": err,
	})
}

func WIPresponder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "method is under construction",
	})
}

// ValidateSignature reports whether signature is a valid HMAC for message.
// https://pkg.go.dev/crypto/hmac
func ValidateSignature(message, signature, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expected := mac.Sum(nil)
	return hmac.Equal(signature, expected)
}
