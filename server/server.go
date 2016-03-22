package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
    
    "io"
    "net/http"
    "golang.org/x/net/websocket"
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
	newClients := make(chan Client)
	deadClients := make(chan Client)
	messages := make(chan Message)

	port := LoadArg("-p", "3000")
	server := launchServer(port)
    launchServerHero("3001")
	go accept(server, newClients)

	for {
		select {
		case client := <-newClients:
			log.Printf("Accepted new client: %v", client.String())
			go handleClient(client, messages, deadClients)
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

func launchServerHero(port string) {
    http.Handle("/echo", websocket.Handler(heroHandler))
    http.Handle("/", http.FileServer(http.Dir(".")))
    err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("ListenAndServe hero: " + err.Error())
	}
    log.Printf("Launching server Hero on %s", port)
}

func heroHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
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
		client, err := authenticate(conn)
		if err == nil {
			clients <- client
		} else {
			s := fmt.Sprintf("Can't authenticate %v, reason : %v", conn.RemoteAddr().String(), err.Error())
			log.Println(s)
			fmt.Fprintf(conn, s)
		}
	}
}

func authenticate(conn net.Conn) (Client, error) {
	reader := bufio.NewReader(conn)
	req := Request{}
	if err := req.Decode(reader); err != nil {
		return Client{}, fmt.Errorf("Fail decode Request from %s", conn.RemoteAddr().String())
	}
	auth := AuthRequest{}
	if err := auth.Decode(&req); err != nil {
		return Client{}, fmt.Errorf("Fail decode AuthRequest from %s", conn.RemoteAddr().String())
	}
	return Client{Conn: conn, ID: auth.ID, Name: auth.Name}, nil
}

func handleClient(client Client, messages chan Message, deadClients chan Client) {
	reader := bufio.NewReader(client.Conn)
	for {
		req := Request{}
		err := req.Decode(reader)
		if err != nil {
			break
		}
		messages <- Message{Client: client, Request: req}
	}
	deadClients <- client
}

func handleMessage(mess Message) {
	log.Printf("From %v : [%v] ", mess.Client.Conn.RemoteAddr(), mess.Request)
}
