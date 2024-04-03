package controllers

import (
	"encoding/json"
	"errors"
	"gym-api/src/auth"
	"gym-api/src/database"
	"gym-api/src/models"
	"gym-api/src/repository"
	"gym-api/src/responses"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateWorkout(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var workout models.Workout
	if err = json.Unmarshal(bodyRequest, &workout); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	workout.UserID = userID

	if err = workout.Prepare(); err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewWorkoutsRepository(db)
	workout.ID, err = repo.Create(workout)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusCreated, workout)
}

func GetWorkoutsByUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	userIDInToken, err := auth.ExtractUserID(r)
	if err != nil {
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIDInToken {
		err = errors.New("it is not possible to view another user's workouts")
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewWorkoutsRepository(db)
	workouts, err := repo.GetWorkoutsNamesByUserID(userID)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusOK, workouts)
}

func UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	workoutID, err := strconv.ParseUint(params["workoutId"], 10, 64)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewWorkoutsRepository(db)
	workoutInDatabase, err := repo.GetWorkoutNameByID(workoutID)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if workoutInDatabase.UserID != userID {
		err = errors.New("it is not possible to update an workout that is not yours")
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var workout models.Workout
	if err = json.Unmarshal(bodyRequest, &workout); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	if err = workout.Prepare(); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	if err = repo.Update(workoutID, userID, workout); err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}

func DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	workoutID, err := strconv.ParseUint(params["workoutId"], 10, 64)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewWorkoutsRepository(db)
	workoutInDatabase, err := repo.GetWorkoutNameByID(workoutID)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if workoutInDatabase.UserID != userID {
		err = errors.New("it is not possible to delete an workout that is not yours")
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if err = repo.Delete(workoutID, userID); err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}

func AddExercises(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	userIDInToken, err := auth.ExtractUserID(r)
	if err != nil {
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIDInToken {
		err = errors.New("it is not possible to add exercises to another user's workout")
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var workout models.Workout
	if err = json.Unmarshal(bodyRequest, &workout); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repoExercise := repository.NewExercisesRepository(db)
	for _, exercise := range workout.Exercises {
		exerciseInDatabase, err := repoExercise.GetExerciseByID(exercise.ID, userID)
		if err != nil {
			responses.SendError(w, http.StatusInternalServerError, err)
			return
		}

		if exerciseInDatabase.UserID != userID {
			err = errors.New("it is not possible to add an exercise to the workout that does not belong to you")
			responses.SendError(w, http.StatusInternalServerError, err)
			return
		}
	}

	repo := repository.NewWorkoutsRepository(db)
	for _, exercise := range workout.Exercises {
		if err = repo.AddExerciseToWorkout(workout.ID, exercise.ID); err != nil {
			responses.SendError(w, http.StatusInternalServerError, err)
			return
		}
	}

	responses.SendJSON(w, http.StatusCreated, nil)
}
