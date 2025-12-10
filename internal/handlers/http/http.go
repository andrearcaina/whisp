package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

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
	KlipyKey string
	router   *gin.Engine // this is set when NewRouter is called. It isn't exported because it should only be used in this pkg
}

func NewHandler(db *db.Database, hub *ws.Hub, klipyKey string) *Handler {
	return &Handler{
		DB:       db,
		Hub:      hub,
		KlipyKey: klipyKey,
	}
}

func (h *Handler) NewRouter() *gin.Engine {
	h.router = gin.New()

	h.router.Use(
		gin.Recovery(),
		middleware.LoggerMiddleware(),
	)

	h.router.Static("/static", "./static")

	h.router.GET("/", h.serveWeb)
	h.router.GET("/ws", h.serveWs)

	h.router.GET("/api/messages", h.listMessages)
	h.router.GET("/api/klipy/gifs/:search", h.listKlipyGifs)
	h.router.GET("/api/klipy/gifs/:search/:limit", h.listKlipyGifs) // limit is optional

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
	log.Println("Fetching messages from the database...")

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
			Message:   &msg.Message.String,
			ImageUrl:  &msg.ImageUrl.String,
			GifUrl:    &msg.GifUrl.String,
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
				"image_url": null,
				"gif_url": null,
				"username": "anonymous",
				"created_at": "2025-09-05T22:08:32.311568Z"
			},
			...
		]
	*/
	c.JSON(http.StatusOK, response)
}

func (h *Handler) listKlipyGifs(c *gin.Context) {
	search := c.Param("search")
	if search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search parameter is required"})
		return
	}

	limit := c.Param("limit")
	if limit == "" {
		limit = "20"
	}

	// properly encode the search parameter to handle spaces and special characters
	endpoint := "https://api.klipy.com/v2/search?key=" + h.KlipyKey + "&q=" + url.QueryEscape(search) + "&limit=" + limit

	if search == "trending" {
		endpoint = "https://api.klipy.com/v2/featured?key=" + h.KlipyKey + "&limit=" + limit
	}
	if search == "stickers" {
		endpoint += "&searchfilter=sticker"
	}

	gifs, err := http.Get(endpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch gifs"})
		return
	}
	if gifs.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch gifs from klipy"})
		return
	}
	defer gifs.Body.Close()

	var klipyResponse klipyGIFResponse
	if err := json.NewDecoder(gifs.Body).Decode(&klipyResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse gifs response"})
		return
	}

	// trim results to only what we need
	// (the gif url technically but also the id and desc for reference as well as putting the desc in the alt attribute of the img tag)
	var response []klipyGIFTrimmedResponse
	for _, result := range klipyResponse.Results {
		response = append(response, klipyGIFTrimmedResponse{
			ID:     result.ID,
			Title:  result.Title,
			GifUrl: result.MediaFormats.GIF.URL,
		})
	}

	/*
		Expected response format:
		[
			{
				"id": "5544765380983011",
				"title": "Bubu the Cat Happy Dancing",
				"gif_url": "https://static.klipy.com/ii/35ccce3d852f7995dd2da910f2abd795/b9/4d/pZ5cIOPM.gif"
			  },
			{
				"id": "2525964843568523",
				"title": "Shaq's Excited Dance",
				"gif_url": "https://static.klipy.com/ii/4e7bea9f7a3371424e6c16ebc93252fe/80/7c/faEUNiNuq5TjUsuXzv.gif"
			  },
			...
		]
	*/
	c.JSON(http.StatusOK, response)
}
