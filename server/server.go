package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Message struct {
	Client Client
	Data   []byte
}

type Client struct {
	Conn net.Conn
	Name string
	Id   string
}

func (c *Client) ToString() string {
	return fmt.Sprintf("Client[Id: %s, Name: %s, Address: %s]", c.Id, c.Name, c.Conn.RemoteAddr())
}

func main() {
	clients := make(map[net.Conn]Client)
	newConnections := make(chan net.Conn)
	deadClients := make(chan Client)
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
			clients[conn] = Client{Conn: conn}
			go read(clients[conn], messages, deadClients)
		case client := <-deadClients:
			log.Printf("%v disconnected", client.ToString())
			delete(clients, client.Conn)
		case message := <-messages:
			go handle(message)
		}
	}
}

func accept(server net.Listener, newConnections chan net.Conn) {
	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		newConnections <- conn
	}
}

func read(client Client, messages chan Message, deadClients chan Client) {
	reader := bufio.NewReader(client.Conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		messages <- Message{client, []byte(incoming)[0 : len(incoming)-1]}
	}
	deadClients <- client
}

func handle(mess Message) {
	log.Printf("From %v : [%v] ", mess.Client.Conn.RemoteAddr(), string(mess.Data))
}
