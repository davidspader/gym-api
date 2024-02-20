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
	{
		URI:            "/users/{userId}",
		Method:         http.MethodPut,
		Function:       controllers.UpdateUser,
		Authentication: true,
	},
	{
		URI:            "/users/{userId}",
		Method:         http.MethodDelete,
		Function:       controllers.DeleteUser,
		Authentication: true,
	},
}
