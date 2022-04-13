package api

import "github.com/gin-gonic/gin"

func abortWithError(c *gin.Context, httpCode int, err string) {
	c.JSON(httpCode, gin.H{
		"error": err,
	})
}
