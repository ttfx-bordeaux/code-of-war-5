package game

import (
	"testing"
)

func TestCreateGame(t *testing.T) {
	game, err := NewGame("name")
	if err != nil || game.ID == "" || game.Name != "name" {
		t.Fail()
	}
}

func TestJoinGame(t *testing.T) {
	game, _ := NewGame("name")
	game.Join(Client{ID: "1", Name: "name1"})
	game.Join(Client{ID: "2", Name: "name2"})
	game.Join(Client{ID: "3", Name: "name3"})
	game.Join(Client{ID: "3", Name: "duplicateID"})

	if len(game.Players) != 3 || game.Players["3"].Name != "name3" {
		t.Fail()
	}
}
