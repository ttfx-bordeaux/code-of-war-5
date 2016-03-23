package io

import (
	"bufio"
	"encoding/json"
)

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

// Message from Client
type Command struct {
	Value string `json:"value"`
}

// Decode Request from bufio.Reader
func (r *Request) Decode(reader *bufio.Reader) error {
	d := json.NewDecoder(reader)
	return d.Decode(r)
}

// Decode AuthRequest from Request
func (a *AuthRequest) Decode(req *Request) error {
	return json.Unmarshal(req.Data, &a)
}

// Decode AuthRequest from Request
func (c *Command) Decode(req *Request) error {
	return json.Unmarshal(req.Data, &c)
}
