package ws

import "log"

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	clients    map[*Client]bool // Registered clients connected to the "hub" (the server)
	broadcast  chan []byte      // Inbound messages from the clients
	register   chan *Client     // Register requests from the clients
	unregister chan *Client     // Unregister requests from clients
}

// NewHub initializes a new Hub instance
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the hub to listen for incoming register, unregister and broadcast requests
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client registered. Total clients: %d", len(h.clients))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			log.Printf("Client unregistered. Total clients: %d", len(h.clients))
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			log.Printf("Broadcasting: %s to %d clients", message, len(h.clients))
		}
	}
}
