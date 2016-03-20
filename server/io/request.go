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

// Decode Request from bufio.Reader
func (r *Request) Decode(reader *bufio.Reader) error {
	incoming, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	err = json.Unmarshal(incoming, &r)
	if err != nil {
		return err
	}
	return nil
}

// Decode AuthRequest from Request
func (a *AuthRequest) Decode(req *Request) (err error) {
	err = json.Unmarshal(req.Data, &a)
	return
}
