package hero

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// http://41j.com/blog/2014/12/gorilla-websockets-golang-simple-websockets-example/
//LaunchServerHero launch websocket server for hero management
func LaunchServerHero(port string) {
	http.Handle("/", http.FileServer(http.Dir("./hero/static")))
	http.Handle("/hero", websocket.Handler(heroHandler))
	http.Handle("/game/inprogress", websocket.Handler(gameProgress))

	log.Printf("Launching server Hero on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func heroHandler(ws *websocket.Conn) {
	var err error

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			log.Printf("Can't receive")
			break
		}

		log.Printf("Received back from client: " + reply)

		msg := "Received:  " + reply
		log.Printf("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			log.Printf("Can't send")
			break
		}
	}
}

func gameProgress(ws *websocket.Conn) {

}
