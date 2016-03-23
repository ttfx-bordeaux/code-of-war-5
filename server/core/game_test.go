package core

import (
	"testing"
)

func TestCreateGame(t *testing.T) {
	clients := make(map[string]Client)
	clients["1"] = Client{ID: "1", Name: "name1"}
	clients["2"] = Client{ID: "2", Name: "name2"}
	clients["3"] = Client{ID: "3", Name: "name3"}
	game, err := NewGame(clients)

	if err != nil || game.ID == "" || len(game.Players) != 3 || game.Players["1"].Name != "name1" {
		t.Fail()
	}
}
