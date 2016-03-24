package core

import (
	"encoding/json"
	"log"

	uuid "github.com/nu7hatch/gouuid"
)

// Game struct
type Game struct {
	ID        string
	Players   map[string]Client
	gameTurns chan (GameTurn)
	Scenario  Scenario
}

// GameTurn represent an action perform  by some entity (client, ia, hero)
type GameTurn struct {
	Type string `json:"type"`
}

// NewGame ceate new game with clients
func NewGame(clients map[string]Client) (Game, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return Game{}, err
	}
	return Game{
		gameTurns: make(chan GameTurn),
		ID:        u4.String(),
		Players:   clients,
		Scenario:  Scenario{},
	}, nil
}

func (g *Game) Launch() {
	log.Printf("Launch Game %s with %d players", g.ID, len(g.Players))
	go g.Scenario.Process(g.gameTurns)
	for _, c := range g.Players {
		log.Printf("Player in da Game %v ", c)
		if err := sendStartMessage(c); err != nil {
			continue
		}
		go c.Process(g.gameTurns)
	}
	go gameTurnsHandler(g.gameTurns)
}

func sendStartMessage(c Client) error {
	e := json.NewEncoder(c.Conn)
	return e.Encode("{insert start mesage here}")
}

func gameTurnsHandler(gts chan GameTurn) {
	for {
		gt := <-gts
		log.Printf("Game turn recived : %s", gt.Type)
	}
}
