package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow Svelte to connect
}

type Message struct {
	Type  string   `json:"type"`  // "update"
	Items []string `json:"items"` // The full list
}

type BasketRoom struct {
	clients map[*websocket.Conn]bool
	items   []string
	mu      sync.Mutex
}

var rooms = make(map[string]*BasketRoom)
var roomsMu sync.Mutex

func main() {
	// 2. WebSocket endpoint
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Backend running on :8080")
	http.ListenAndServe(":8080", nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	basketID := r.URL.Query().Get("id")
	if basketID == "" {
		return
	}

	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	// LOCK the rooms map to check/create safely
	roomsMu.Lock()
	room, exists := rooms[basketID]
	if !exists {
		// Lazy Initialization: Create the basket because someone arrived
		room = &BasketRoom{
			clients: make(map[*websocket.Conn]bool),
			items:   []string{},
		}
		rooms[basketID] = room
		fmt.Printf("Created new basket: %s\n", basketID)
	}
	roomsMu.Unlock()

	// Register the user to the room
	room.mu.Lock()
	room.clients[conn] = true
	// Send existing items (will be empty list if new)
	conn.WriteJSON(Message{Type: "update", Items: room.items})
	room.mu.Unlock()

	// Listen for updates...
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			room.mu.Lock()
			delete(room.clients, conn) // Clean up on disconnect
			room.mu.Unlock()
			break
		}

		room.mu.Lock()
		room.items = msg.Items
		for client := range room.clients {
			client.WriteJSON(msg)
		}
		room.mu.Unlock()
	}
}
