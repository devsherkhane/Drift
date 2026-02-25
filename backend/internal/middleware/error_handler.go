package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler catches any panics in the request chain and returns a JSON error
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error detail for debugging
				log.Printf("Panic recovered: %v", err)

				// Return a clean JSON response instead of crashing
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "An unexpected internal server error occurred",
				})
			}
		}()
		c.Next()
	}
}