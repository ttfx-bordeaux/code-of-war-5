package main

import (
	"errors"
	"net"
	"testing"
	"time"
)

func TestAccept(t *testing.T) {
	done := make(chan bool)
	connections := make(chan net.Conn)

	var conn net.Conn
	go func() {
		conn = <-connections
		done <- true
	}()
	go accept(&stubListener{}, connections)

	<-done
	if conn.RemoteAddr().String() != "1.2.3.4:5" {
		t.Fail()
	}
}

func TestDontAccept(t *testing.T) {
	done := make(chan bool)
	connections := make(chan net.Conn)

	go func() {
		go accept(&stubListener{fail: true}, connections)
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

type stubListener struct {
	fail bool
}

func (m stubListener) Accept() (net.Conn, error) {
	if m.fail {
		return nil, errors.New("fail")
	}
	return stubConn{}, nil
}
func (m stubListener) Close() error   { return nil }
func (m stubListener) Addr() net.Addr { return nil }

type stubConn struct {
}

func (m stubConn) RemoteAddr() net.Addr               { return stubAddr{} }
func (m stubConn) Write(b []byte) (n int, err error)  { return 0, nil }
func (m stubConn) SetDeadline(t time.Time) error      { return nil }
func (m stubConn) SetReadDeadline(t time.Time) error  { return nil }
func (m stubConn) SetWriteDeadline(t time.Time) error { return nil }
func (m stubConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (m stubConn) LocalAddr() net.Addr                { return nil }
func (m stubConn) Close() error                       { return nil }

type stubAddr struct {
}

func (m stubAddr) Network() string { return "network" }
func (m stubAddr) String() string  { return "1.2.3.4:5" }
