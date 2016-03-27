package io

import (
	"log"
	"net"
)

// Server handle net.Listener and encapsulate handler
type Server interface {
	net.Listener
}

// Accepter : Accept connection
type Accepter interface {
	Accept() (net.Conn, error)
}

// LaunchServer create a server on port with his handler
func LaunchServer(port string, handler func(Accepter)) Server {
	lst, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	log.Printf("Launching server on %s", lst.Addr())
	go handler(lst)
	return lst
}
