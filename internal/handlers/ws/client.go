package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/andrearcaina/whisp/internal/db"
	"github.com/andrearcaina/whisp/internal/db/generated"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgtype"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

type IncomingMessage struct {
	Message string `json:"message"`
	GifUrl  string `json:"gif_url,omitempty"`
}

type OutgoingMessage struct {
	ID        int32     `json:"id"`
	Message   string    `json:"message,omitempty"`
	GifUrl    string    `json:"gif_url,omitempty"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Client) readPump(db *db.Database) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		// Parse the incoming JSON message
		var incomingMsg IncomingMessage
		if err := json.Unmarshal(message, &incomingMsg); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		params := generated.CreateMessageParams{
			Message: pgtype.Text{String: incomingMsg.Message, Valid: incomingMsg.Message != ""},
			GifUrl:  pgtype.Text{String: incomingMsg.GifUrl, Valid: incomingMsg.GifUrl != ""},
		}

		// Store the message and gif_url in the database
		resp, err := db.GetQueries().CreateMessage(context.Background(), params)
		if err != nil {
			log.Printf("Failed to save message: %v", err)
			continue
		}

		// Create outgoing message with proper structure
		outgoingMsg := OutgoingMessage{
			ID:        resp.ID,
			Message:   resp.Message.String,
			GifUrl:    resp.GifUrl.String,
			Username:  "anonymous",
			CreatedAt: resp.CreatedAt.Time,
		}

		// Decode the outgoing message to JSON
		msgBytes, err := json.Marshal(outgoingMsg)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			continue
		}

		/*
			Broadcast the message to all clients
			Expected response format:
			{
				"id": 33,
				"message": "🐟",
				"image_url": null,
				"gif_url": null,
				"username": "anonymous",
				"created_at": "2025-09-05T22:08:32.311568Z"
			  }
		*/
		c.hub.broadcast <- msgBytes
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, db *db.Database, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump(db)
}
