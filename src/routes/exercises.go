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
	{
		URI:            "/exercises/{exerciseId}",
		Method:         http.MethodGet,
		Function:       controllers.GetExercise,
		Authentication: true,
	},
}
