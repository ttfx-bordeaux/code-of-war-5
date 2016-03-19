package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/nu7hatch/gouuid"
)

// Message from Client
type Message struct {
	Client Client
	Data   []byte
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
	newConnections := make(chan net.Conn)
	deadClients := make(chan Client)
	messages := make(chan Message)

	port := loadArg("-p", "3000")
	server := launchServer(port)
	go accept(server, newConnections)

	for {
		select {
		case conn := <-newConnections:
			if uuid, err := uuid.NewV4(); err == nil {
				c := Client{ID: uuid.String(), Conn: conn}
				log.Printf("Accepted new client: %v", c.String())
				go read(c, messages, deadClients)
			}
			log.Printf("Can't create uuid for Client from : %v", conn.RemoteAddr().String())
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

func accept(server Accepter, newConnections chan net.Conn) {
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
