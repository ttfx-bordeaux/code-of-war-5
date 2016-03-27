package io

import (
	"log"
	"net"
)

// Server handle net.Listener and encapsulate handler
type Server struct {
	listener net.Listener
}

// Accepter : Accept connection
type Accepter interface {
	Accept() (net.Conn, error)
}

func LaunchServer(port string, handler func(Accepter)) *Server {
	lst, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	log.Printf("Launching server on %s", lst.Addr())
	go handler(lst)
	return &Server{listener: lst}
}

func (s *Server) Close() error {
	log.Printf("closing server on %s", s.Addr())
	return s.listener.Close()
}

func (s *Server) Addr() string {
	return s.listener.Addr().String()
}
