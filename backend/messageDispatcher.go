package main

import (
	"encoding/json"
)

type MethodHandler func(client *Client, payload json.RawMessage) error

// MethodDispatcher maps method strings to their respective handler functions.
type MethodDispatcher struct {
	handlers map[string]MethodHandler
}

// NewMethodDispatcher creates a new dispatcher with an initialized map.
func NewMethodDispatcher() *MethodDispatcher {
	return &MethodDispatcher{
		handlers: make(map[string]MethodHandler),
	}
}

func (d *MethodDispatcher) Register(methodName string, handler MethodHandler) {
	d.handlers[methodName] = handler
}

func (d *MethodDispatcher) GetHandler(methodName string) (MethodHandler, bool) {
	handler, ok := d.handlers[methodName]
	return handler, ok
}
