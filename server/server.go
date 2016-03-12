package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Message struct {
	Connection net.Conn
	Data       []byte
}

func main() {
	clients := make(map[net.Conn]string)
	newConnections := make(chan net.Conn)
	deadConnections := make(chan net.Conn)
	messages := make(chan Message)

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
			go read(conn, messages, deadConnections)
		case conn := <-deadConnections:
			log.Printf("Client %v disconnected", clients[conn])
			delete(clients, conn)
		case message := <-messages:
			go handle(message)
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

func read(conn net.Conn, messages chan Message, deadConnections chan net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		messages <- Message{conn, []byte(incoming)[0 : len(incoming)-1]}
	}
	deadConnections <- conn
}

func handle(mess Message) {
	log.Printf("From %v : [%v] ", mess.Connection.RemoteAddr(), string(mess.Data))
}
