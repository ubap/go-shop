package main

import "encoding/json"

type Message struct {
	MessageId string          `json:"messageId"`
	Method    string          `json:"method"`
	Payload   json.RawMessage `json:"payload"`
}
