package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	URI            string
	Method         string
	Function       func(http.ResponseWriter, *http.Request)
	Authentication bool
}

func configRoutes(r *mux.Router) *mux.Router {
	routes := usersRoutes

	for _, route := range routes {
		r.HandleFunc(route.URI, route.Function).Methods(route.Method)
	}

	return r
}

func GenerateRouter() *mux.Router {
	r := mux.NewRouter()
	return configRoutes(r)
}
