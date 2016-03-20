package core

import (
	"fmt"
	"net"
)

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
