package routes

import (
	"gym-api/src/controllers"
	"net/http"
)

type UserControllers struct {
	UserControllers *controllers.UserController
}

var usersRoutes = []route{
	{
		URI:            "/users",
		Method:         http.MethodPost,
		Function:       controllers.UserController.CreateUser,
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
	{
		URI:            "/users/{userId}/change-password",
		Method:         http.MethodPost,
		Function:       controllers.ChangePassword,
		Authentication: true,
	},
}
