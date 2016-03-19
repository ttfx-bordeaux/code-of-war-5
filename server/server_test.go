package main

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"testing"
)

func TestAccept(t *testing.T) {
	done := make(chan bool)
	connections := make(chan net.Conn)

	var conn net.Conn
	go func() {
		conn = <-connections
		done <- true
	}()
	go accept(&accepterPass{}, connections)

	<-done
	if conn.RemoteAddr().String() != "1.2.3.4:5" {
		t.Fail()
	}
}

func TestDontAccept(t *testing.T) {
	done := make(chan bool)
	connections := make(chan net.Conn)

	go func() {
		go accept(&accepterFail{}, connections)
		done <- true
	}()

	<-done
	if len(connections) > 0 {
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

func TestReadMessageFromClient(t *testing.T) {
	done := make(chan bool)
	deadClients := make(chan Client)
	messages := make(chan Message)
	mockConn := newStubConn()
	var msg Message

	go func() {
		msg = <-messages
		done <- true
	}()
	go read(Client{Conn: mockConn}, messages, deadClients)
	auth := &AuthRequest{ID: "12345", Name: "kriyss"}
	b, _ := json.Marshal(auth)
	req := Request{Action: "authenticate", Data: b}
	encoder := json.NewEncoder(mockConn.ClientWriter)
	encoder.Encode(&req)

	<-done
	if string(msg.Request.Action) != "authenticate" {
		t.Fail()
	}
}

type accepterPass struct {
	Accepter
}

func (m accepterPass) Accept() (net.Conn, error) {
	return stubConn{}, nil
}

type accepterFail struct {
	Accepter
}

func (m accepterFail) Accept() (net.Conn, error) {
	return nil, errors.New("fail")
}

type stubConn struct {
	net.Conn
	ServerReader *io.PipeReader
	ServerWriter *io.PipeWriter
	ClientReader *io.PipeReader
	ClientWriter *io.PipeWriter
}

func newStubConn() stubConn {
	serverRead, clientWrite := io.Pipe()
	clientRead, serverWrite := io.Pipe()
	return stubConn{
		ServerReader: serverRead,
		ServerWriter: serverWrite,
		ClientReader: clientRead,
		ClientWriter: clientWrite,
	}
}

func (m stubConn) Read(data []byte) (n int, err error)  { return m.ServerReader.Read(data) }
func (m stubConn) Write(data []byte) (n int, err error) { return m.ServerWriter.Write(data) }
func (m stubConn) RemoteAddr() net.Addr                 { return stubAddr{} }

type stubAddr struct {
}

func (m stubAddr) Network() string { return "network" }
func (m stubAddr) String() string  { return "1.2.3.4:5" }
