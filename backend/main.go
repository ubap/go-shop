package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"go-shop/backend/db"
	"io/fs"
	"log"
	"net/http"
	"strings"
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

//go:embed "dist/*"
var frontendAssets embed.FS

var rooms = make(map[string]*BasketRoom)
var roomsMu sync.Mutex

var store db.Store

func main() {
	sqliteStore, err := db.NewSqliteStore("db.sqlite")
	if err != nil {
		log.Fatal(err)
		return
	}
	store = sqliteStore

	http.Handle("/", SPAHandler())
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Running backend on :9090")
	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println("Closing")
}

func SPAHandler() http.HandlerFunc {
	distFS, err := fs.Sub(frontendAssets, "dist")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(distFS))

	return func(w http.ResponseWriter, r *http.Request) {
		originalPath := r.URL.Path
		cleanPath := strings.TrimPrefix(originalPath, "/")
		isFile := false
		if cleanPath != "" {
			if _, err := fs.Stat(distFS, cleanPath); err == nil {
				isFile = true
			}
		}
		if isFile && originalPath != "/" && originalPath != "/index.html" {
			fileServer.ServeHTTP(w, r)
			return
		}
		indexBytes, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		var initialData []db.Item = nil
		if strings.HasPrefix(originalPath, "/basket/") {
			basketKey := strings.TrimPrefix(originalPath, "/basket/")
			initialData, err = store.GetItemsForBasket(basketKey)
		}
		jsonData, err := json.Marshal(initialData)
		if err != nil {
			jsonData = []byte("null")
		}
		placeholder := []byte(`null /* INJECT_INITIAL_DATA */`)
		finalHTML := bytes.Replace(indexBytes, placeholder, jsonData, 1)
		w.Write(finalHTML)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	basketID := r.URL.Query().Get("id")
	if basketID == "" {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	defer func(conn *websocket.Conn) {
		fmt.Printf("Closing connection: %v\n", conn.RemoteAddr())
		err := conn.Close()
		if err != nil {
			fmt.Printf("Error closing websocket connection: %v\n", conn.RemoteAddr())
		}
	}(conn)
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

		// TODO Extract from here
		trimmed := strings.TrimSpace(msg.ItemName)
		runes := []rune(trimmed)
		// 2. Limit to 55 characters
		if len(runes) > 55 {
			msg.ItemName = string(runes[:55])
		} else {
			msg.ItemName = trimmed
		}

		switch msg.Type {
		case C2SAddItem:
			room.AddItem(msg.ItemName)
		case C2SDeleteItem:
			room.DeleteItem(msg.Id)
		case C2SRestoreItem:
			room.RestoreItem(msg.Id)
		case C2SSetItemCompletion:
			room.SetItemCompletion(msg.Id, msg.Completed)
		}

		fmt.Printf("Received message for basket %s: %+v\n", basketID, msg)
	}
}
