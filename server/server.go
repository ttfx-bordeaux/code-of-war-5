package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/websocket"
	"net/http"

	"github.com/ttfx-bordeaux/code-of-war-5/server/core"
	"github.com/ttfx-bordeaux/code-of-war-5/server/io"
	"github.com/ttfx-bordeaux/code-of-war-5/server/util"
)

// Message from Client
type Message struct {
	Client  core.Client
	Request io.Request
}

var (
	// ConnectedClients : all authentified clients
	ConnectedClients map[string]core.Client
)

func main() {
	ConnectedClients = make(map[string]core.Client)
	newClients := make(chan core.Client)
	deadClients := make(chan core.Client)
	messages := make(chan Message)
	commands := make(chan io.Command)

	gamePort := util.LoadArg("-p", "3000")
	gameServer := launchServer(gamePort)
	go accept(gameServer, newClients)

	commandPort := util.LoadArg("-p", "4000")
	commandServer := launchServer(commandPort)
	go acceptCommand(commandServer, commands)

	launchServerHero("3001")

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
		case c := <-commands:
			log.Printf("Command received %+v", c)
			if c.Value == "start" {
				game, err := core.NewGame(ConnectedClients)
				if err != nil {
					log.Println(err)
					continue
				}
				game.Launch()

			}
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
	// package io already used
	//io.Copy(ws, ws)
}

// Accepter : Accept connection
type Accepter interface {
	Accept() (net.Conn, error)
}

func accept(server Accepter, clients chan core.Client) {
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

func acceptCommand(server Accepter, commands chan io.Command) {
	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()
		r := bufio.NewReader(conn)
		req := io.Request{}
		req.Decode(r)
		cmd := io.Command{}
		cmd.Decode(&req)
		commands <- cmd
	}
}

// DuplicateClientIDErr : error when multiple client use same ID
type DuplicateClientIDErr struct {
	error
}

func (err DuplicateClientIDErr) Error() string {
	return "ID already in use"
}

func authenticate(conn net.Conn, connected map[string]core.Client) (core.Client, error) {
	r := bufio.NewReader(conn)
	req := io.Request{}
	if err := req.Decode(r); err != nil {
		return core.Client{}, err
	}
	auth := io.AuthRequest{}
	if err := auth.Decode(&req); err != nil {
		return core.Client{}, err
	}
	if _, exist := connected[auth.ID]; exist {
		return core.Client{}, DuplicateClientIDErr{}
	}
	return core.Client{Conn: conn, ID: auth.ID, Name: auth.Name}, nil
}

func handleClient(client core.Client, messages chan Message, deadClients chan core.Client) {
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
