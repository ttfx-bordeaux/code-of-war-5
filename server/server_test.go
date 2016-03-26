package main

import "testing"

func TestLaunchServer(t *testing.T) {
	srv := LaunchServer("2000", func(a Accepter) {})
	if srv.Addr() != "[::]:2000" {
		t.Fail()
	}
}

func TestFailLaunchServer(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fail()
		}
	}()
	LaunchServer("-1", func(new Accepter) { panic("error") })
}

func TestCloseServer(t *testing.T) {
	srv := LaunchServer("2001", func(a Accepter) {})
	srv.Close()
}

func TestHandler(t *testing.T) {
	quit := make(chan bool)
	LaunchServer("2001", func(a Accepter) { quit <- true })
	if b := <-quit; !b {
		t.Fail()
	}
}
