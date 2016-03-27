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

func TestQuitGame(t *testing.T) {
	c1 := Client{ID: "1", Name: "name1"}
	c3 := Client{ID: "3", Name: "name3"}

	game, _ := NewGame("name")
	game.Join(c1)
	game.Join(Client{ID: "2", Name: "name2"})
	game.Join(c3)

	game.Remove(c1)
	game.Remove(c3)

	if len(game.Players) != 1 || game.Players["2"].Name != "name2" {
		t.Fail()
	}
}
