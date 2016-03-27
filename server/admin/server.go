package admin

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"time"

	"github.com/gorilla/mux"
)

//Route define path
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes define array of paths
type Routes []Route

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		index,
	},
	Route{
		"Game",
		"POST",
		"/game",
		game,
	},
	Route{
		"GameShow",
		"POST",
		"/game/{gameId}/launch",
		gameLaunch,
	},
}

//LaunchServerAdmin launch a server for administration purpose
func LaunchServerAdmin(port string) {
	router := newRouter()
	log.Printf("Launching server Admin on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func game(w http.ResponseWriter, r *http.Request) {
	gameID := "1233hhjh22"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(gameID); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Create Game:", gameID)
}

func gameLaunch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameId"]
	fmt.Fprintln(w, "LaunchGame:", gameID)
}

//Logger log
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
