package http

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/andrearcaina/whisp/internal/db"
	"github.com/andrearcaina/whisp/internal/db/generated"
	"github.com/andrearcaina/whisp/internal/handlers/ws"
	"github.com/andrearcaina/whisp/views"
	"github.com/gin-gonic/gin"
)

func NewRouter(db *db.Database, hub *ws.Hub) *gin.Engine {
	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		serveWeb(c)
	})

	r.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, db, c.Writer, c.Request)
	})

	r.GET("/api/messages", func(c *gin.Context) {
		listMessages(db, c)
	})

	return r
}

func serveWeb(c *gin.Context) {
	component := views.ChatPage()
	templ.Handler(component).ServeHTTP(c.Writer, c.Request)
}

func listMessages(db *db.Database, c *gin.Context) {
	messages, err := db.GetQueries().ListMessages(c.Request.Context(), generated.ListMessagesParams{
		Limit:  50,
		Offset: 0,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
		return
	}

	var response []messageResponse
	for _, msg := range messages {
		response = append(response, messageResponse{
			ID:        msg.ID,
			Message:   msg.Message,
			Username:  "anonymous",
			CreatedAt: msg.CreatedAt.Time,
		})
	}

	c.JSON(http.StatusOK, response)
}
