package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") || path == "/" || path == "/ws" {
			log.Printf("%s %s", c.Request.Method, path)
		}
		c.Next()
	}
}
