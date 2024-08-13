package routes

import (
	"gym-api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
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
	routes = append(routes, exercisesRoutes...)
	routes = append(routes, workoutsRotues...)

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

	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs/"))))
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	))

	return r
}

func GenerateRouter() *mux.Router {
	r := mux.NewRouter()
	return configRoutes(r)
}
