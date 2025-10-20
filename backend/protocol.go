package main

import (
	"encoding/json"
	"fmt"
	"go-shop/backend/basket"
	"go-shop/backend/basket/inmemory"
	"log"

	"github.com/google/uuid"
)

type Protocol struct {
	basket   *inmemory.Basket
	handlers map[string]MethodHandler
}

func NewProtocol(basket *inmemory.Basket) *Protocol {
	p := &Protocol{basket, make(map[string]MethodHandler)}
	p.handlers["itemUpdate"] = p.onItemUpdate
	return p
}

func (p *Protocol) onConnect(c *Client) {
	items := (*p.basket).GetAllItems()

	for _, item := range items {
		marshal, err := json.Marshal(item)
		if err != nil {
			fmt.Errorf("error marshaling item: %v", err)
		}

		message, err := json.Marshal(Message{
			MessageId: uuid.New().String(), Method: "itemUpdate", Payload: marshal})
		if err != nil {
			fmt.Errorf("error marshaling message: %v", err)
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

		marshal, err := json.Marshal(Message{MessageId: msg.MessageId, Method: "ack"})
		if err != nil {
			log.Printf("error marshalling ack: %v", err)
			return
		}
		client.send <- marshal

		// TODO move broadcast to specific method handler,
		// it depends on the method whether we want to handle it
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
