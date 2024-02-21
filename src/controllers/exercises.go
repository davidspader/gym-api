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

func CreateExercise(w http.ResponseWriter, r *http.Request) {
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

	var exercise models.Exercise
	if err = json.Unmarshal(bodyRequest, &exercise); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	exercise.UserID = userID

	if err = exercise.Prepare(); err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewExercisesRepository(db)
	exercise.ID, err = repo.Create(exercise)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusCreated, exercise)
}

func GetExercise(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	exerciseID, err := strconv.ParseUint(params["exerciseId"], 10, 64)
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

	repo := repository.NewExercisesRepository(db)
	exercise, err := repo.GetExerciseByID(exerciseID, userID)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusOK, exercise)
}

func GetExercisesByUser(w http.ResponseWriter, r *http.Request) {
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
		err = errors.New("it is not possible to view another user's exercises")
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewExercisesRepository(db)
	exercises, err := repo.GetExercisesByUserID(userID)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusOK, exercises)
}
