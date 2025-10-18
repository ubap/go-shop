package main

import (
	"encoding/json"
	"fmt"
	"go-shop/backend/basket"
	"log"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	dispatcher *MethodDispatcher
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		handler, ok := c.dispatcher.GetHandler(msg.Method)
		if ok {
			handler(c, msg.Payload)

			log.Printf("Received valid method '%s', broadcasting...", msg.Method)
			broadcastMsg := &BroadcastMessage{
				Sender:  c,
				Message: messageBytes,
			}
			c.hub.broadcast <- broadcastMsg
		} else {
			log.Printf("Received unknown method: %s", msg.Method)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			// The hub closed the channel.
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		_, err = w.Write(message)
		if err != nil {
			return
		}

		if err := w.Close(); err != nil {
			return
		}
	}
}

func (c *Client) onConnect() {
	items := basket.GetAllItems()

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

func handleItemUpdate(client *Client, payload json.RawMessage) error {
	log.Printf("Handling 'handleItemUpdate' from client %p with payload: %s", client, string(payload))

	var basketItem basket.BasketItem
	if err := json.Unmarshal(payload, &basketItem); err != nil {
		log.Printf("error unmarshaling message: %v", err)
	}

	basket.UpdateItem(basketItem)

	return nil // Return nil if successful
}
