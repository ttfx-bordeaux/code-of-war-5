package core

import (
	"fmt"
	"net"
	"time"
)

type Render interface {
	Process(gameTurns chan GameTurn)
}

// Client connected
type Client struct {
	net.Conn
	Name string
	ID   string
}

type Scenario struct {
	Name string
}

// String : format client information
func (c *Client) String() string {
	return fmt.Sprintf("Client[Id: %s, Name: %s, Address: %s]", c.ID, c.Name, c.Conn.RemoteAddr())
}

func (c Client) Process(gameTurns chan GameTurn) {
	for {
		time.Sleep(2 * time.Second)
		gameTurns <- GameTurn{c.ID + "_ACTION"}
	}
}

func (s Scenario) Process(gameTurns chan GameTurn) {
	for {
		time.Sleep(2 * time.Second)
		gameTurns <- GameTurn{s.Name + "_ACTION_SCENARIO"}
	}
}
