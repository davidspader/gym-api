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
	{
		URI:            "/users/{userId}/exercises",
		Method:         http.MethodGet,
		Function:       controllers.GetExercisesByUser,
		Authentication: true,
	},
	{
		URI:            "/exercises/{exerciseId}",
		Method:         http.MethodPut,
		Function:       controllers.UpdateExercise,
		Authentication: true,
	},
	{
		URI:            "/exercises/{exerciseId}",
		Method:         http.MethodDelete,
		Function:       controllers.DeleteExercise,
		Authentication: true,
	},
}
