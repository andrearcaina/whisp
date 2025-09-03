package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/andrearcaina/whisp/internal/db"
	"github.com/gorilla/websocket"
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
}

type OutgoingMessage struct {
	ID        int32     `json:"id"`
	Message   string    `json:"message"`
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

		// Store only the content in the database
		resp, err := db.GetQueries().CreateMessage(context.Background(), incomingMsg.Message)
		if err != nil {
			log.Printf("Failed to save message: %v", err)
			continue
		}

		// Create outgoing message with proper structure
		outgoingMsg := OutgoingMessage{
			ID:        resp.ID,
			Message:   resp.Message,
			Username:  "anonymous",
			CreatedAt: resp.CreatedAt.Time,
		}

		// Marshal to JSON and broadcast to all clients
		msgBytes, err := json.Marshal(outgoingMsg)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			continue
		}

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
