package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"

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

	gamePort := util.LoadArg("-p", "3000")
	gameSrv := LaunchServer(gamePort, clientHandler)
	defer gameSrv.Close()

	commandPort := util.LoadArg("-p", "4000")
	adminSrv := LaunchServer(commandPort, commandHandler)
	defer adminSrv.Close()

	for {
	}
}

func clientHandler(server Accepter) {
	newClients := make(chan core.Client)
	deadClients := make(chan core.Client)
	go handleNewConnections(server, newClients)
	go handleClients(newClients, deadClients)
	go handleDisconnectedClient(deadClients)
}

func handleDisconnectedClient(deadClients chan core.Client) {
	for {
		c := <-deadClients
		delete(ConnectedClients, c.ID)
		log.Printf("%v disconnected", c.String())
	}
}

func handleClients(clients, deadClients chan core.Client) {
	for {
		c := <-clients
		log.Printf("Accepted new client: %v", c.String())
		ConnectedClients[c.ID] = c
		go func(client core.Client, deadClients chan core.Client) {
			r := bufio.NewReader(client.Conn)
			for {
				req := io.Request{}
				if err := req.Decode(r); err != nil {
					break
				}
			}
			deadClients <- client
		}(c, deadClients)
	}
}
func handleNewConnections(server Accepter, clients chan core.Client) {
	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()
		err = authenticate(conn, ConnectedClients, func(c core.Client) { clients <- c })
		if err != nil {
			s := fmt.Sprintf("Can't authenticate %v, reason : %v", conn.RemoteAddr().String(), err.Error())
			log.Println(s)
			fmt.Fprintf(conn, s)
		}
	}
}

func commandHandler(server Accepter) {
	commands := make(chan io.Command)
	go func(cmds chan io.Command) {
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
			cmds <- cmd
		}
	}(commands)

	go func(cmds chan io.Command) {
		for {
			c := <-cmds
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
	}(commands)
}

func authenticate(conn net.Conn, connected map[string]core.Client, fct func(client core.Client)) error {
	r := bufio.NewReader(conn)
	req := io.Request{}
	if err := req.Decode(r); err != nil {
		return err
	}
	auth := io.AuthRequest{}
	if err := auth.Decode(&req); err != nil {
		return err
	}
	if _, exist := connected[auth.ID]; exist {
		return errors.New("DuplicateID")
	}
	fct(core.Client{Conn: conn, ID: auth.ID, Name: auth.Name})
	return nil
}
