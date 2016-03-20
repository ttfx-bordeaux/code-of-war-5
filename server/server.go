package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/ttfx-bordeaux/code-of-war-5/server/io"
	"github.com/ttfx-bordeaux/code-of-war-5/server/util"
)

// Message from Client
type Message struct {
	Client  Client
	Request io.Request
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

var (
	// ConnectedClients : all authentified clients
	ConnectedClients map[string]Client
)

func main() {
	ConnectedClients = make(map[string]Client)
	newClients := make(chan Client)
	deadClients := make(chan Client)
	messages := make(chan Message)

	port := util.LoadArg("-p", "3000")
	server := launchServer(port)
	go accept(server, newClients)

	for {
		select {
		case c := <-newClients:
			log.Printf("Accepted new client: %v", c.String())
			ConnectedClients[c.ID] = c
			go handleClient(c, messages, deadClients)
		case c := <-deadClients:
			delete(ConnectedClients, c.ID)
			log.Printf("%v disconnected", c.String())
		case m := <-messages:
			log.Printf("%+v", ConnectedClients)
			go handleMessage(m)
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
		defer conn.Close()
		c, err := authenticate(conn, ConnectedClients)
		if err == nil {
			clients <- c
		} else {
			s := fmt.Sprintf("Can't authenticate %v, reason : %v", conn.RemoteAddr().String(), err.Error())
			log.Println(s)
			fmt.Fprintf(conn, s)
		}
	}
}

// DuplicateClientIDErr : error when multiple client use same ID
type DuplicateClientIDErr struct {
	error
}

func (err DuplicateClientIDErr) Error() string {
	return "ID already in use"
}

func authenticate(conn net.Conn, connected map[string]Client) (Client, error) {
	r := bufio.NewReader(conn)
	req := io.Request{}
	if err := req.Decode(r); err != nil {
		return Client{}, err
	}
	auth := io.AuthRequest{}
	if err := auth.Decode(&req); err != nil {
		return Client{}, err
	}
	if _, exist := connected[auth.ID]; exist {
		return Client{}, DuplicateClientIDErr{}
	}
	return Client{Conn: conn, ID: auth.ID, Name: auth.Name}, nil
}

func handleClient(client Client, messages chan Message, deadClients chan Client) {
	r := bufio.NewReader(client.Conn)
	for {
		req := io.Request{}
		if err := req.Decode(r); err != nil {
			break
		}
		messages <- Message{Client: client, Request: req}
	}
	deadClients <- client
}

func handleMessage(mess Message) {
	log.Printf("From %v : [%v] ", mess.Client.Conn.RemoteAddr(), mess.Request)
}
