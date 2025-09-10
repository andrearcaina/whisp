package http

import (
	"encoding/json"
	"net/http"

	"github.com/andrearcaina/whisp/internal/db"
	"github.com/andrearcaina/whisp/internal/db/generated"
	"github.com/andrearcaina/whisp/internal/handlers/ws"
	"github.com/andrearcaina/whisp/internal/middleware"
	"github.com/andrearcaina/whisp/views"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB       *db.Database
	Hub      *ws.Hub
	TenorKey string
	router   *gin.Engine // this is set when NewRouter is called. It isn't exported because it should only be used in this pkg
}

func NewHandler(db *db.Database, hub *ws.Hub, tenorKey string) *Handler {
	return &Handler{
		DB:       db,
		Hub:      hub,
		TenorKey: tenorKey,
	}
}

func (h *Handler) NewRouter() *gin.Engine {
	h.router = gin.New()

	h.router.Use(
		gin.Recovery(),
		middleware.BotMiddleware(),
		middleware.LoggerMiddleware(),
	)

	h.router.Static("/static", "./static")

	h.router.GET("/", h.serveWeb)
	h.router.GET("/ws", h.serveWs)

	h.router.GET("/api/messages", h.listMessages)
	h.router.GET("/api/tenor/gifs/:search", h.listTenorGifs)
	h.router.GET("/api/tenor/gifs/:search/:limit", h.listTenorGifs) // limit is optional

	return h.router
}

func (h *Handler) serveWeb(c *gin.Context) {
	if err := views.ChatPage().Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "failed to render page")
		return
	}
}

func (h *Handler) serveWs(c *gin.Context) {
	ws.ServeWs(h.Hub, h.DB, c.Writer, c.Request)
}

func (h *Handler) listMessages(c *gin.Context) {
	messages, err := h.DB.GetQueries().ListMessages(c.Request.Context(), generated.ListMessagesParams{
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

	/*
		Expected response format:
		[
			{
				"id": 33,
				"message": "🐟",
				"username": "anonymous",
				"created_at": "2025-09-05T22:08:32.311568Z"
			},
			...
		]
	*/
	c.JSON(http.StatusOK, response)
}

func (h *Handler) listTenorGifs(c *gin.Context) {
	search := c.Param("search")
	limit := c.Param("limit")
	if search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search parameter is required"})
		return
	}
	if limit == "" {
		limit = "10"
	}

	url := "https://tenor.googleapis.com/v2/search?q=" + search + "&key=" + h.TenorKey + "&client_key=whisp_app" + "&limit=" + limit

	gifs, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch gifs"})
		return
	}
	if gifs.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch gifs from tenor"})
		return
	}
	defer gifs.Body.Close()

	var tenorResponse tenorGIFResponse
	if err := json.NewDecoder(gifs.Body).Decode(&tenorResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse gifs response"})
		return
	}

	// trim results to only what we need
	// (the gif url technically but also the id and desc for reference as well as putting the desc in the alt attribute of the img tag)
	var response []tenorGIFTrimmedResponse
	for _, result := range tenorResponse.Results {
		response = append(response, tenorGIFTrimmedResponse{
			ID:     result.ID,
			Desc:   result.ContentDescription,
			GifUrl: result.MediaFormats.GIF.URL,
		})
	}

	/*
		Expected response format:
		[
			{
				"id": "15784368690949001113",
				"desc": "a man in a blue shirt and tie is screaming with his hands in the air",
				"gif_url": "https://media.tenor.com/2w1XsfvQD5kAAAAC/hhgf.gif"
			  },
			{
				"id": "16963945808433096689",
				"desc": "a little girl is laughing with her fist in the air while wearing a vest and tie .",
				"gif_url": "https://media.tenor.com/62wK1Xyhp_EAAAAC/happy.gif"
			  },
			...
		]
	*/
	c.JSON(http.StatusOK, response)
}
