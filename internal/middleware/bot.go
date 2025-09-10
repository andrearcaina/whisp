package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func BotMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// do not allow access to any path that is not /ws, /api or /static
		if !(path == "/" ||
			path == "/ws" ||
			strings.HasPrefix(path, "/static/") ||
			strings.HasPrefix(path, "/api/messages") ||
			strings.HasPrefix(path, "/api/tenor/gifs")) {
			c.AbortWithStatus(404)
			return
		}

		c.Next()
	}
}
