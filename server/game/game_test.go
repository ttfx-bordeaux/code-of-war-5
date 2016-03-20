package game

import (
	"testing"

	"github.com/ttfx-bordeaux/code-of-war-5/server/core"
)

func TestCreateGame(t *testing.T) {
	clients := make(map[string]core.Client)
	clients["1"] = core.Client{ID: "1", Name: "name1"}
	clients["2"] = core.Client{ID: "2", Name: "name2"}
	clients["3"] = core.Client{ID: "3", Name: "name3"}
	game := NewGame(clients)

	if game.ID == "" || len(game.Players) != 3 || game.Players["1"].Name != "name1" {
		t.Fail()
	}
}
