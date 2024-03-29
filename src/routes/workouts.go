package routes

import (
	"gym-api/src/controllers"
	"net/http"
)

var workoutsRotues = []route{
	{
		URI:            "/workouts",
		Method:         http.MethodPost,
		Function:       controllers.CreateWorkout,
		Authentication: true,
	},
	{
		URI:            "/workouts/user/{userId}",
		Method:         http.MethodGet,
		Function:       controllers.GetWorkoutsByUser,
		Authentication: true,
	},
}
