package main

import (
	"errors"
	"net"
	"testing"
)

func TestDontAccept(t *testing.T) {
	done := make(chan bool)
	clients := make(chan Client)

	go func() {
		go accept(&accepterFail{}, clients)
		done <- true
	}()

	<-done
	if len(clients) > 0 {
		t.Fail()
	}
}

func TestLaunchServer(t *testing.T) {
	srv := launchServer("2000")
	if srv.Addr().String() != "[::]:2000" {
		t.Fail()
	}
}

func TestFailLaunchServer(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fail()
		}
	}()
	launchServer("-1")
}

type accepterFail struct {
	Accepter
}

func (m accepterFail) Accept() (net.Conn, error) {
	return nil, errors.New("fail")
}
