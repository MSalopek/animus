package engine

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AbortErr(c *gin.Context, httpCode int, err error) {
	c.Abort()
	c.JSON(httpCode, gin.H{
		"error": err.Error(),
	})
}

func WIPresponder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "method is under construction",
	})
}

func LogRequest(logger *log.Logger) gin.HandlerFunc {
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
