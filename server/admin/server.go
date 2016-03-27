package admin

import (
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

func newRouter(routes Routes) *mux.Router {
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

//LaunchServerAdmin launch a server for administration purpose
func LaunchServerAdmin(port string, routes Routes) {
	log.Printf("Launching server Admin on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, newRouter(routes)))
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
