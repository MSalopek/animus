package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: remove
func WIPresponder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "method is under construction",
	})
}

func abortWithError(c *gin.Context, httpCode int, err string) {
	c.JSON(httpCode, gin.H{
		"error": err,
	})
}
