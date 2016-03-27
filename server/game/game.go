package game

import (
	"encoding/json"
	"log"

	uuid "github.com/nu7hatch/gouuid"
)

// Game struct
type Game struct {
	ID        string
	Name      string
	Players   map[string]Client
	gameTurns chan (GameTurn)
	Scenario  Scenario
}

// GameTurn represent an action perform  by some entity (client, ia, hero)
type GameTurn struct {
	Type string `json:"type"`
}

// NewGame ceate new game with clients
func NewGame(name string) (Game, error) {
	u4, err := uuid.NewV4()
	if err != nil {
		return Game{}, err
	}
	log.Printf("create game [%s:%s]", name, u4.String())
	return Game{
		ID:      u4.String(),
		Players: make(map[string]Client),
		Name:    name,
	}, nil
}

// Join a client to Game
func (g *Game) Join(c Client) {
	if _, exist := g.Players[c.ID]; !exist {
		g.Players[c.ID] = c
		log.Printf("add client [%s:%s] to game [%s:%s]", c.Name, c.ID, g.Name, g.ID)
	}
}

// Remove client from the Game
func (g *Game) Remove(c Client) {
	if _, exist := g.Players[c.ID]; exist {
		delete(g.Players, c.ID)
		log.Printf("remove client [%s:%s] to game [%s:%s]", c.Name, c.ID, g.Name, g.ID)
	}
}

func (g *Game) Launch() {
	log.Printf("launch game [%s:%s] with %d players", g.Name, g.ID, len(g.Players))
	go g.Scenario.Process(g.gameTurns)
	for _, c := range g.Players {
		log.Printf("Player in da Game %s ", c.String())
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
		log.Printf("Game turn received : %s", gt.Type)
	}
}
