package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ttfx-bordeaux/code-of-war-5/server/admin"
	"github.com/ttfx-bordeaux/code-of-war-5/server/game"
	"github.com/ttfx-bordeaux/code-of-war-5/server/hero"
	"github.com/ttfx-bordeaux/code-of-war-5/server/io"
	"github.com/ttfx-bordeaux/code-of-war-5/server/util"
)

var (
	// ConnectedClients : all authentified clients
	ConnectedClients map[string]game.Client

	// AdminActions : actions that admin can do
	AdminActions map[string]func()

	// AllGame that are created
	AllGame map[string]game.Game

	routes = admin.Routes{
		admin.Route{
			"Index",
			"GET",
			"/",
			index,
		},
		admin.Route{
			"Game",
			"POST",
			"/game",
			createGame,
		},
		admin.Route{
			"GameJoin",
			"POST",
			"/game/{gameId}/player/{playerId}",
			joinGame,
		},
		admin.Route{
			"GameShow",
			"POST",
			"/game/{gameId}/launch",
			launchGame,
		},
	}
)

func main() {
	ConnectedClients = make(map[string]game.Client)
	AllGame = make(map[string]game.Game)

	gamePort := util.LoadArg("--port", "3000")
	gameSrv := io.LaunchServer(gamePort, game.NewHandler(ConnectedClients))
	defer gameSrv.Close()

	adminPort := util.LoadArg("--admin-port", "4000")
	go admin.LaunchServerAdmin(adminPort, routes)

	heroPort := util.LoadArg("--hero-port", "4001")
	go hero.LaunchServerHero(heroPort)

	for {
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func createGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	g, err := game.NewGame(vars["gameName"])
	if err != nil {
		log.Println(err)
	}
	AllGame[g.ID] = g
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(g.ID); err != nil {
		panic(err)
	}
}

func launchGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	g := AllGame[vars["gameId"]]
	g.Launch()
}

func joinGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	g := AllGame[vars["gameId"]]
	p := ConnectedClients[vars["playerId"]]
	g.Join(p)

}
