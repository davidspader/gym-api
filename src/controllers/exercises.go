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

	var exercise models.Exercise
	if err = json.Unmarshal(bodyRequest, &exercise); err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	exercise.UserID = userID

	if err = exercise.Prepare(); err != nil {
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

	repo := repository.NewExercisesRepository(db)
	exercise.ID, err = repo.Save(exercise)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusCreated, exercise)
}

func GetExercise(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		err = errors.New(responses.ErrMsgUnauthorized)
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	exerciseID, err := strconv.ParseUint(params["exerciseId"], 10, 64)
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

	repo := repository.NewExercisesRepository(db)
	exercise, err := repo.FindByID(exerciseID, userID)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if exercise.ID == 0 {
		err = errors.New("the exercise you are looking for could not be found")
		responses.SendError(w, http.StatusNotFound, err)
		return
	}

	responses.SendJSON(w, http.StatusOK, exercise)
}

func GetExercisesByUser(w http.ResponseWriter, r *http.Request) {
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
		err = errors.New("it is not possible to view another user's exercises")
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

	repo := repository.NewExercisesRepository(db)
	exercises, err := repo.FindByUserID(userID)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusOK, exercises)
}

func UpdateExercise(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		err = errors.New(responses.ErrMsgUnauthorized)
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	exerciseID, err := strconv.ParseUint(params["exerciseId"], 10, 64)
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

	repo := repository.NewExercisesRepository(db)
	exerciseInDatabase, err := repo.FindByID(exerciseID, userID)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if exerciseInDatabase.UserID != userID {
		err = errors.New("it is not possible to update an exercise that is not yours")
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.New(responses.ErrMsgUnprocessableEntity)
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var exercise models.Exercise
	if err = json.Unmarshal(bodyRequest, &exercise); err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	if exercise.Name == "" {
		exercise.Name = exerciseInDatabase.Name
	}

	if exercise.Weight == 0 {
		exercise.Weight = exerciseInDatabase.Weight
	}

	if exercise.Reps == 0 {
		exercise.Reps = exerciseInDatabase.Reps
	}

	if err = exercise.Prepare(); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	if err = repo.Update(exerciseID, userID, exercise); err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}

func DeleteExercise(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ExtractUserID(r)
	if err != nil {
		err = errors.New(responses.ErrMsgUnauthorized)
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	params := mux.Vars(r)
	exerciseID, err := strconv.ParseUint(params["exerciseId"], 10, 64)
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

	repo := repository.NewExercisesRepository(db)
	exerciseInDatabase, err := repo.FindByID(exerciseID, userID)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if exerciseInDatabase.UserID != userID {
		err = errors.New("it is not possible to delete an exercise that is not yours")
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	if err = repo.Delete(exerciseID, userID); err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}
