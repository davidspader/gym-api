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
		err = errors.New(responses.ErrMsgUnauthorized)
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.New(responses.ErrMsgUnprocessableEntity)
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var workout models.Workout
	if err = json.Unmarshal(bodyRequest, &workout); err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
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
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewWorkoutsRepository(db)
	workout.ID, err = repo.Save(workout)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusCreated, workout)
}

func GetWorkoutsByUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	userIDInToken, err := auth.ExtractUserID(r)
	if err != nil {
		err = errors.New(responses.ErrMsgUnauthorized)
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
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewWorkoutsRepository(db)
	workouts, err := repo.FindNamesByUserID(userID)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusOK, workouts)
}

func UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		err = errors.New(responses.ErrMsgUnauthorized)
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	workoutID, err := strconv.ParseUint(params["workoutId"], 10, 64)
	if err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewWorkoutsRepository(db)
	workoutInDatabase, err := repo.FindNameByID(workoutID)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if workoutInDatabase.UserID != userID {
		err = errors.New("it is not possible to update an workout that is not yours")
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.New(responses.ErrMsgUnprocessableEntity)
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var workout models.Workout
	if err = json.Unmarshal(bodyRequest, &workout); err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	if err = workout.Prepare(); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	if err = repo.Update(workoutID, userID, workout); err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}

func DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		err = errors.New(responses.ErrMsgUnprocessableEntity)
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	workoutID, err := strconv.ParseUint(params["workoutId"], 10, 64)
	if err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewWorkoutsRepository(db)
	workoutInDatabase, err := repo.FindNameByID(workoutID)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if workoutInDatabase.UserID != userID {
		err = errors.New("it is not possible to delete an workout that is not yours")
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	if err = repo.Delete(workoutID, userID); err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}

func GetWorkout(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		err = errors.New(responses.ErrMsgUnauthorized)
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	workoutID, err := strconv.ParseUint(params["workoutId"], 10, 64)
	if err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewWorkoutsRepository(db)
	workout, err := repo.FindByID(workoutID, userID)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if workout.ID == 0 {
		err = errors.New("the workout you are looking for could not be found")
		responses.SendError(w, http.StatusNotFound, err)
		return
	}

	responses.SendJSON(w, http.StatusOK, workout)
}

func AddExercises(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	workoutID, err := strconv.ParseUint(params["workoutId"], 10, 64)
	if err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := auth.ExtractUserID(r)
	if err != nil {
		err = errors.New(responses.ErrMsgUnauthorized)
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.New(responses.ErrMsgUnprocessableEntity)
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var exerciseIDs models.ExerciseIDs
	if err = json.Unmarshal(bodyRequest, &exerciseIDs); err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repoExercise := repository.NewExercisesRepository(db)
	if err = repoExercise.VerifyOwnership(exerciseIDs.IDs, userID); err != nil {
		err = errors.New(responses.ErrMsgForbidden)
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	repo := repository.NewWorkoutsRepository(db)
	if err = repo.AddExercises(workoutID, exerciseIDs.IDs); err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}

func RemoveExercises(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	workoutID, err := strconv.ParseUint(params["workoutId"], 10, 64)
	if err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := auth.ExtractUserID(r)
	if err != nil {
		err = errors.New(responses.ErrMsgUnauthorized)
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.New(responses.ErrMsgUnprocessableEntity)
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var exerciseIDs models.ExerciseIDs
	if err = json.Unmarshal(bodyRequest, &exerciseIDs); err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repoExercise := repository.NewExercisesRepository(db)
	if err = repoExercise.VerifyOwnership(exerciseIDs.IDs, userID); err != nil {
		err = errors.New(responses.ErrMsgForbidden)
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	repo := repository.NewWorkoutsRepository(db)
	if err = repo.RemoveExercises(workoutID, exerciseIDs.IDs); err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}
