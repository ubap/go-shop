package main

import (
	"encoding/json"
	"log"
)

type BasketItem struct {
	ItemID       string          `json:"id"`
	Name         string          `json:"name"`
	LastModified json.RawMessage `json:"lastModified"`
}

type Basket interface {
}

// this method belongs to Client I think
func handleItemUpdate(client *Client, payload json.RawMessage) error {
	log.Printf("Handling 'handleItemUpdate' from client %p with payload: %s", client, string(payload))

	var basketItem BasketItem
	if err := json.Unmarshal(payload, &basketItem); err != nil {
		log.Printf("error unmarshaling message: %v", err)
	}

	UpdateItem(basketItem)

	return nil // Return nil if successful
}

// global variable
var basketItems = make(map[string]BasketItem)

func UpdateItem(basketItem BasketItem) {
	basketItems[basketItem.ItemID] = basketItem
}

func GetAllItems() []BasketItem {
	basketItemsSlice := make([]BasketItem, 0, len(basketItems))

	for _, value := range basketItems {
		basketItemsSlice = append(basketItemsSlice, value)
	}
	return basketItemsSlice
}
