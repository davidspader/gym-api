package routes

import (
	"gym-api/src/controllers"
	"net/http"
)

var usersRoutes = []route{
	{
		URI:            "/users",
		Method:         http.MethodPost,
		Function:       controllers.CreateUser,
		Authentication: false,
	},
}
