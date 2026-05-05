package main

import (
	"fmt"
	"go-shop/backend/db"
	"sync"

	"github.com/gorilla/websocket"
)

type BasketRoom struct {
	basketKey string

	clients map[*websocket.Conn]bool
	store   db.Store

	mu sync.Mutex
}

func (b *BasketRoom) LogInUser(conn *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()

	basketItems, err := b.store.GetItemsForBasket(b.basketKey)
	if err != nil {
		return
	}

	// register client
	b.clients[conn] = true

	// send initial data
	conn.WriteJSON(S2CMessage{Type: S2CFullList, Items: basketItems})
}

func (b *BasketRoom) Remove(conn *websocket.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.clients, conn)
}

func (b *BasketRoom) Broadcast(msg S2CMessage) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for client := range b.clients {
		err := client.WriteJSON(msg)

		if err != nil {
			delete(b.clients, client)
		}
	}
}

func (b *BasketRoom) UpdateClients() {
	basketItems, err := b.store.GetItemsForBasket(b.basketKey)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	b.Broadcast(S2CMessage{Type: S2CFullList, Items: basketItems})
}

func (b *BasketRoom) AddItem(itemName string) {
	b.store.AddItemToBasket(b.basketKey, itemName)
	b.UpdateClients()
}

func (b *BasketRoom) DeleteItem(itemId int64) {
	b.store.DeleteItem(b.basketKey, itemId)
	b.UpdateClients()
}

func (b *BasketRoom) SetItemCompletion(itemId int64, completed bool) {
	store.SetItemCompletion(b.basketKey, itemId, completed)
	b.UpdateClients()
}
