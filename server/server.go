package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	clients := make(map[net.Conn]string)
	newConnections := make(chan net.Conn)
	deadConnections := make(chan net.Conn)

	server, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Launching server on %s", server.Addr())

	go accept(server, newConnections)

	for {
		select {
		case conn := <-newConnections:
			addr := conn.RemoteAddr().String()
			log.Printf("Accepted new client, %v", addr)
			clients[conn] = addr
			go read(conn, deadConnections)
		case conn := <-deadConnections:
			log.Printf("Client %v disconnected", clients[conn])
			delete(clients, conn)
		}
	}
}

func accept(server net.Listener, newConnections chan net.Conn) {
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
		}
		newConnections <- conn
	}
}

func read(conn net.Conn, deadConnections chan net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		log.Printf("Received : %v", incoming)
	}
	deadConnections <- conn
}
