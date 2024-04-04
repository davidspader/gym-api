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
		URI:            "/users/{userId}/workouts",
		Method:         http.MethodGet,
		Function:       controllers.GetWorkoutsByUser,
		Authentication: true,
	},
	{
		URI:            "/workouts/{workoutId}",
		Method:         http.MethodPut,
		Function:       controllers.UpdateWorkout,
		Authentication: true,
	},
	{
		URI:            "/workouts/{workoutId}",
		Method:         http.MethodDelete,
		Function:       controllers.DeleteWorkout,
		Authentication: true,
	},
	{
		URI:            "/workouts/{workoutId}",
		Method:         http.MethodGet,
		Function:       controllers.GetWorkout,
		Authentication: true,
	},
	{
		URI:            "/workouts/addExercise/{userId}",
		Method:         http.MethodPost,
		Function:       controllers.AddExercises,
		Authentication: true,
	},
	{
		URI:            "/workouts/removeExercise/{userId}",
		Method:         http.MethodDelete,
		Function:       controllers.RemoveExercises,
		Authentication: true,
	},
}
