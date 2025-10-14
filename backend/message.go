package main

import "encoding/json"

type Message struct {
	Method  string          `json:"method"`
	Payload json.RawMessage `json:"payload"`
}
