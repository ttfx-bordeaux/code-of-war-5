package game

import (
	"encoding/json"
	"log"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/ttfx-bordeaux/code-of-war-5/server/core"
)

// Game struct
type Game struct {
	ID      string
	Players map[string]core.Client
}

// NewGame ceate new game with clients
func NewGame(clients map[string]core.Client) (Game, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return Game{}, err
	}
	return Game{
		ID:      u4.String(),
		Players: clients,
	}, nil
}

func (g *Game) Launch() {
	log.Printf("Launch Game %s with %d players", g.ID, len(g.Players))
	for _, c := range g.Players {
		e := json.NewEncoder(c.Conn)
		e.Encode(g)
	}
}
