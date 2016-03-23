package core

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Render interface {
	Process()
}

// Client connected
type Client struct {
	net.Conn
	Name string
	ID   string
}

// String : format client information
func (c *Client) String() string {
	return fmt.Sprintf("Client[Id: %s, Name: %s, Address: %s]", c.ID, c.Name, c.Conn.RemoteAddr())
}

func (c Client) Process(gameTurns chan GameTurn) {
	log.Println("process_" + c.ID)
	for {
		time.Sleep(2 * time.Second)
		gameTurns <- GameTurn{c.ID + "_ACTION"}
	}
}
