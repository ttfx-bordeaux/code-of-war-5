package main

import (
	"log"

	"github.com/ttfx-bordeaux/code-of-war-5/server/admin"
	"github.com/ttfx-bordeaux/code-of-war-5/server/game"
	"github.com/ttfx-bordeaux/code-of-war-5/server/io"
	"github.com/ttfx-bordeaux/code-of-war-5/server/util"
)

// Message from Client
type Message struct {
	Client  game.Client
	Request io.Request
}

var (
	// ConnectedClients : all authentified clients
	ConnectedClients map[string]game.Client

	// AdminActions : actions that admin can do
	AdminActions map[string]func()

	// AllGame that are created
	AllGame map[string]game.Game
)

func main() {
	ConnectedClients = make(map[string]game.Client)
	AllGame = make(map[string]game.Game)
	initAdminActions()

	gamePort := util.LoadArg("--port", "3000")
	gameSrv := io.LaunchServer(gamePort, game.NewHandler(ConnectedClients))
	defer gameSrv.Close()

	commandPort := util.LoadArg("--admin-port", "4000")
	adminSrv := io.LaunchServer(commandPort, admin.NewHandler(AdminActions))
	defer adminSrv.Close()

	for {
	}
}

func initAdminActions() {
	AdminActions = map[string]func(){
		"create": func() {
			g, err := game.NewGame(ConnectedClients)
			if err != nil {
				log.Println(err)
			}
			// g.Launch()
			AllGame[g.ID] = g
		},
		"all-game": func() {
			log.Printf("%+v", AllGame)
		},
		"all-player": func() {
			log.Printf("%+v", ConnectedClients)
		},
	}
}
