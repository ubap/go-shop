package main

import (
	"encoding/json"
	"log"
)

type BasketItem struct {
	Name         string          `json:"name"`
	LastModified json.RawMessage `json:"lastModified"`
}

// handleItemAdded is the handler for the "itemAddedToBuy" method.
func handleItemAdded(client *Client, payload json.RawMessage) error {
	log.Printf("Handling 'itemAddedToBuy' from client %p with payload: %s", client, string(payload))

	var basketItem BasketItem
	if err := json.Unmarshal(payload, &basketItem); err != nil {
		log.Printf("error unmarshaling message: %v", err)
	}

	return nil // Return nil if successful
}

// handleItemBought is the handler for the "itemBought" method.
func handleItemBought(client *Client, payload json.RawMessage) error {
	log.Printf("Handling 'itemBought' from client %p with payload: %s", client, string(payload))
	return nil
}
