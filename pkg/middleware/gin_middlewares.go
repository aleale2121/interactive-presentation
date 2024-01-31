package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ContentTypeChecker is a middleware to check content type.
func ContentTypeChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if contentType != "application/json" {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Specified content type not allowed."})
			c.Abort()
			return
		}
		c.Next()
	}
}
