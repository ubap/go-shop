package main

import (
	"encoding/json"
	"fmt"
	"go-shop/backend/basket"
	"log"
)

type Protocol struct {
	basket   *basket.InMemoryBasket
	handlers map[string]MethodHandler
}

func NewProtocol(basket *basket.InMemoryBasket) *Protocol {
	p := &Protocol{basket, make(map[string]MethodHandler)}
	p.handlers["itemUpdate"] = p.onItemUpdate
	return p
}

func (p *Protocol) onConnect(c *Client) {
	items := (*p.basket).GetAllItems()

	for _, item := range items {
		marshal, err := json.Marshal(item)
		if err != nil {
			fmt.Errorf("Error marshaling item: %v", err)
		}

		message, err := json.Marshal(Message{"itemUpdate", marshal})
		if err != nil {
			fmt.Errorf("Error marshaling message: %v", err)
		}

		c.send <- message
	}

}

func (p *Protocol) onItemUpdate(client *Client, payload json.RawMessage) error {
	log.Printf("Handling 'onItemUpdate' from client %p with payload: %s", client, string(payload))

	var basketItem basket.Item
	if err := json.Unmarshal(payload, &basketItem); err != nil {
		log.Printf("error unmarshaling message: %v", err)
	}

	(*p.basket).UpsertItem(basketItem)

	return nil // Return nil if successful
}

func (p *Protocol) handleMsg(client *Client, messageBytes []byte) {
	var msg Message
	if err := json.Unmarshal(messageBytes, &msg); err != nil {
		log.Printf("error unmarshaling message: %v", err)
	}

	handler, ok := p.getHandler(msg.Method)
	if ok {
		err := handler(client, msg.Payload)
		if err != nil {
			log.Printf("error handling message: %v", err)
			return
		}

		log.Printf("Received valid method '%s', broadcasting...", msg.Method)
		broadcastMsg := &BroadcastMessage{
			Sender:  client,
			Message: messageBytes,
		}
		client.hub.broadcast <- broadcastMsg
	} else {
		log.Printf("Received unknown method: %s", msg.Method)
	}
}

type MethodHandler func(client *Client, payload json.RawMessage) error

func (p *Protocol) getHandler(methodName string) (MethodHandler, bool) {
	handler, ok := p.handlers[methodName]
	return handler, ok
}
