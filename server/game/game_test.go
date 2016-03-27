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
