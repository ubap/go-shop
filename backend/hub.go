package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	// Registered clients. We use a map with a boolean value for quick additions and deletions.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *BroadcastMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Mutex to protect access to the clients map
	mu sync.Mutex

	dispatcher *MethodDispatcher
}

// BroadcastMessage is a wrapper for a message that includes the sender,
// so we can avoid sending the message back to its originator.
type BroadcastMessage struct {
	Sender  *Client
	Message []byte
}

// NewHub creates a new Hub instance.
func NewHub(dispatcher *MethodDispatcher) *Hub {
	return &Hub{
		broadcast:  make(chan *BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		dispatcher: dispatcher,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Println("Client registered")
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("Client unregistered")
			}
			h.mu.Unlock()
		case broadcastMsg := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				// Don't send the message back to the sender
				if client != broadcastMsg.Sender {
					select {
					case client.send <- broadcastMsg.Message:
					default:
						// If the send buffer is full, close the client.
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// upgrader upgrades HTTP connections to WebSocket connections.
// We set a generous buffer size and disable origin checking for development.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In production, you should check the origin of the request.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), dispatcher: hub.dispatcher}
	client.hub.register <- client
	client.onConnect()

	go client.writePump()
	go client.readPump()
}
