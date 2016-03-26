package hero

import(
	"log"
    "io"
	"net/http"
    
    "golang.org/x/net/websocket"
)

//LaunchServerHero launch websocket server for hero management
func LaunchServerHero(port string) {
	http.Handle("/", http.FileServer(http.Dir("./hero/static")))
    http.Handle("/hero", websocket.Handler(heroHandler))
	
    log.Printf("Launching server Hero on %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func heroHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}