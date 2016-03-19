package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

// Message from Client
type Message struct {
	Client  Client
	Request Request
}

// Client connected
type Client struct {
	Conn net.Conn
	Name string
	ID   string
}

// String : format client information
func (c *Client) String() string {
	return fmt.Sprintf("Client[Id: %s, Name: %s, Address: %s]", c.ID, c.Name, c.Conn.RemoteAddr())
}

func main() {
	connectedClients := []Client{}
	newClients := make(chan Client)
	deadClients := make(chan Client)
	messages := make(chan Message)

	port := loadArg("-p", "3000")
	server := launchServer(port)
	go accept(server, newClients)

	for {
		select {
		case client := <-newClients:
			log.Printf("Accepted new client: %v", client.String())
			connectedClients = append(connectedClients, client)
			go read(client, messages, deadClients)
		case client := <-deadClients:
			log.Printf("%v disconnected", client.String())
		case message := <-messages:
			go handleMessage(message)
		}
	}
}

func launchServer(port string) net.Listener {
	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	log.Printf("Launching server on %s", server.Addr())
	return server
}

// Accepter : Accept connection
type Accepter interface {
	Accept() (net.Conn, error)
}

func accept(server Accepter, clients chan Client) {
	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		clients <- Client{Conn: conn}
	}
}

func read(client Client, messages chan Message, deadClients chan Client) {
	reader := bufio.NewReader(client.Conn)
	for {
		incoming, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		var req Request
		err = json.Unmarshal(incoming, &req)
		if err != nil {
			log.Printf("For client %s, can't parse request: %s", client.String(), string(incoming))
		}
		messages <- Message{client, req}
	}
	deadClients <- client
}

func handleMessage(mess Message) {
	log.Printf("From %v : [%v] ", mess.Client.Conn.RemoteAddr(), mess.Request)
}

func loadArg(command, defaultValue string) (port string) {
	port = defaultValue
	args := os.Args
	for i, arg := range args {
		switch {
		case arg == command:
			port = args[i+1]
		}
	}
	return
}
