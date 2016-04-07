package game

import "testing"

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
	game, _ := NewGame("name")
	game.Join(Client{ID: "1", Name: "name1"})
	game.Join(Client{ID: "2", Name: "name2"})
	game.Join(Client{ID: "3", Name: "name3"})

	game.Remove(Client{ID: "1"})
	game.Remove(Client{ID: "3"})

	if len(game.Players) != 1 || game.Players["2"].Name != "name2" {
		t.Fail()
	}
}

func TestLaunchGame(t *testing.T) {
	game, _ := NewGame("my game")
	game.Join(Client{ID: "1", Name: "name1"})
	game.Join(Client{ID: "2", Name: "name2"})
	game.Join(Client{ID: "3", Name: "name3"})

	game.Launch()

	for _, p := range game.Players {
		if len(game.Maps[p.ID]) != 200 || len(game.Maps[p.ID][0]) != 20 {
			t.Fatalf("Map dimension aren't correct for player %s", p.String())
		}
	}
}
