package routes

import (
	"gym-api/src/controllers"
	"net/http"
)

var rotaLogin = route{
	URI:            "/login",
	Method:         http.MethodPost,
	Function:       controllers.Login,
	Authentication: false,
}
