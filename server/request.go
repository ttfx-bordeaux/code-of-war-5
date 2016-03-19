package main

import "encoding/json"

// Request : main structure to dialog with server
type Request struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

// AuthRequest : structure for authenticate client
type AuthRequest struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
