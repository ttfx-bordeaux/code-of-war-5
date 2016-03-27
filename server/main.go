package main

import (
	"log"

	"github.com/ttfx-bordeaux/code-of-war-5/server/admin"
	"github.com/ttfx-bordeaux/code-of-war-5/server/game"
	"github.com/ttfx-bordeaux/code-of-war-5/server/hero"
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
	AdminActions     map[string]func()
)

func main() {
	ConnectedClients = make(map[string]game.Client)
	initAdminActions()

	gamePort := util.LoadArg("--port", "3000")
	gameSrv := io.LaunchServer(gamePort, game.NewHandler(ConnectedClients))
	defer gameSrv.Close()

	commandPort := util.LoadArg("--admin-port", "4000")
	adminSrv := io.LaunchServer(commandPort, admin.NewHandler(AdminActions))
	defer adminSrv.Close()

	go admin.LaunchServerAdmin("3002")
	go hero.LaunchServerHero("3001")

	for {
	}
}

func initAdminActions() {
	AdminActions = map[string]func(){
		"start": func() {
			game, err := game.NewGame(ConnectedClients)
			if err != nil {
				log.Println(err)
			}
			game.Launch()
		},
	}
}
