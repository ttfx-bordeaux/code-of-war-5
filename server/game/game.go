package game

import (
	"encoding/json"
	"log"

	"github.com/ttfx-bordeaux/code-of-war-5/server/core"
)

// Game struct
type Game struct {
	ID      string
	Players map[string]core.Client
}

// NewGame ceate new game with clients
func NewGame(clients map[string]core.Client) Game {
	return Game{
		ID:      "id",
		Players: clients,
	}
}

func (g *Game) Launch() {
	log.Printf("Launch Game %s with %d players", g.ID, len(g.Players))
	for _, c := range g.Players {
		e := json.NewEncoder(c.Conn)
		e.Encode("start")
	}
}
