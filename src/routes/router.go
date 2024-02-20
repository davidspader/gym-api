package routes

import (
	"gym-api/src/middlewares"
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
	routes = append(routes, rotaLogin)

	for _, route := range routes {
		if route.Authentication {
			r.HandleFunc(
				route.URI,
				middlewares.Logger(middlewares.VerifyAuthenticateUser(route.Function)),
			).Methods(route.Method)
		} else {
			r.HandleFunc(
				route.URI,
				middlewares.Logger(route.Function),
			).Methods(route.Method)
		}
	}

	return r
}

func GenerateRouter() *mux.Router {
	r := mux.NewRouter()
	return configRoutes(r)
}
