package main

import (
	"fmt"
	"go-shop/backend/db"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow Svelte to connect
}

type C2SMessage struct {
	Type C2SCommand `json:"type"`

	ItemName string `json:"itemName"`

	Id        int64 `json:"id"`
	Completed bool  `json:"completed"`
}

type S2CMessage struct {
	Type  S2CCommand `json:"type"`
	Items []db.Item  `json:"items"`
}

var rooms = make(map[string]*BasketRoom)
var roomsMu sync.Mutex

var store db.Store

func main() {

	sqliteStore, err := db.NewSqliteStore("db.sqlite")
	if err != nil {
		return
	}
	store = sqliteStore

	// 2. WebSocket endpoint
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Running backend on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println("Closing")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	basketID := r.URL.Query().Get("id")
	if basketID == "" {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Upgrade error: %v\n", err)
		return
	}

	// LOCK the rooms map to check/create safely
	roomsMu.Lock()
	room, exists := rooms[basketID]
	if !exists {
		// Lazy Initialization: Create the basket because someone arrived
		room = &BasketRoom{
			basketKey: basketID,
			clients:   make(map[*websocket.Conn]bool),
			store:     store,
		}
		rooms[basketID] = room
		fmt.Printf("Created new basket: %s\n", basketID)
	}
	roomsMu.Unlock()

	room.LogInUser(conn)

	// Listen for updates...
	for {
		var msg C2SMessage
		if err := conn.ReadJSON(&msg); err != nil {
			fmt.Printf("Error reading from client: %s\n", err)
			room.Remove(conn)
			break
		}

		switch msg.Type {
		case C2SAddItem:
			room.AddItem(msg.ItemName)
		case C2SSetItemCompletion:
			room.SetItemCompletion(msg.Id, msg.Completed)
		}

		fmt.Printf("Received message for basket %s: %+v\n", basketID, msg)
	}
}
