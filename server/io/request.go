package io

import (
	"bufio"
	"encoding/json"
	"fmt"
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
		return RequestDecodeErr{err}
	}
	err = json.Unmarshal(incoming, &r)
	if err != nil {
		return AuthRequestDecodeErr{err}
	}
	return nil
}

// Decode AuthRequest from Request
func (a *AuthRequest) Decode(req *Request) (err error) {
	err = json.Unmarshal(req.Data, &a)
	return
}

// DecodeErr : can't parse structure
type DecodeErr interface {
	error
}

// RequestDecodeErr : can't parse Request structure
type RequestDecodeErr struct {
	DecodeErr
}

// AuthRequestDecodeErr : can't parse AuthRequest structure
type AuthRequestDecodeErr struct {
	DecodeErr
}

func (e AuthRequest) Error() string {
	return fmt.Sprintf("Can't parse AuthRequest : %v", e.Error())
}

func (e RequestDecodeErr) Error() string {
	return fmt.Sprintf("Can't parse RequestDecodeErr : %v", e.Error())
}
