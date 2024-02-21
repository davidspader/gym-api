package routes

import (
	"gym-api/src/controllers"
	"net/http"
)

var exercisesRoutes = []route{
	{
		URI:            "/exercises",
		Method:         http.MethodPost,
		Function:       controllers.CreateExercise,
		Authentication: true,
	},
}
