package io

import (
	"bufio"
	"log"
	"strings"
	"testing"
)

func TestDecodeRequest(t *testing.T) {
	s := "{\"action\":\"authenticate\",\"data\":{ \"example\": \"RAW_JSON\"}}\n"
	req := Request{}

	req.Decode(bufio.NewReader(strings.NewReader(s)))

	if req.Action != "authenticate" {
		t.Fail()
	}
}

func TestFailDecodeRequest(t *testing.T) {
	req := Request{}
	r := bufio.NewReader(strings.NewReader("bad parsing"))
	if err := req.Decode(r); err == nil {
		t.Fail()
	}
}

func TestDecodeAuthRequest(t *testing.T) {
	s := "{\"name\":\"kriyss\",\"id\":\"12345\"}"
	req := Request{Action: "authenticate", Data: []byte(s)}
	auth := AuthRequest{}

	auth.Decode(&req)

	if auth.ID != "12345" {
		log.Println("Fail decode AuthRequest.ID")
		t.Fail()
	}
	if auth.Name != "kriyss" {
		log.Println("Fail decode AuthRequest.Name")
		t.Fail()
	}
}
