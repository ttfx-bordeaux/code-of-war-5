package game

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/ttfx-bordeaux/code-of-war-5/server/io"
)

var (
	// ConnectedClients : all authentified clients
	connectedClients map[string]Client
)

func NewHandler(connClients map[string]Client) func(io.Accepter) {
	connectedClients = connClients
	return func(server io.Accepter) {
		newClients := make(chan Client)
		deadClients := make(chan Client)
		go handleNewConnections(server, newClients)
		go handleClients(newClients, deadClients)
		go handleDisconnectedClient(deadClients)
	}
}

func handleDisconnectedClient(deadClients chan Client) {
	for {
		c := <-deadClients
		delete(connectedClients, c.ID)
		log.Printf("%v disconnected", c.String())
	}
}

func handleClients(clients, deadClients chan Client) {
	for {
		c := <-clients
		log.Printf("Accepted new client: %v", c.String())
		connectedClients[c.ID] = c
		go func(client Client, deadClients chan Client) {
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
func handleNewConnections(server io.Accepter, clients chan Client) {
	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()
		err = authenticate(conn, connectedClients, func(c Client) { clients <- c })
		if err != nil {
			s := fmt.Sprintf("Can't authenticate %v, reason : %v", conn.RemoteAddr().String(), err.Error())
			log.Println(s)
			fmt.Fprintf(conn, s)
		}
	}
}

func authenticate(conn net.Conn, connected map[string]Client, fct func(client Client)) error {
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
	fct(Client{Conn: conn, ID: auth.ID, Name: auth.Name})
	return nil
}
