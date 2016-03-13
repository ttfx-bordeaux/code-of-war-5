package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
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
	newConnections := make(chan net.Conn)
	deadClients := make(chan Client)
	messages := make(chan Message)

	port := loadArg("-p", "3000")
	server := launchServer(port)
	go accept(server, newConnections)

	for {
		select {
		case conn := <-newConnections:
			log.Printf("Accepted new client, %v", conn.RemoteAddr().String())
			go read(Client{Conn: conn}, messages, deadClients)
		case client := <-deadClients:
			log.Printf("%v disconnected", client.ToString())
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

func handleMessage(mess Message) {
	log.Printf("From %v : [%v] ", mess.Client.Conn.RemoteAddr(), string(mess.Data))
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
