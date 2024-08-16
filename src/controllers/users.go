package controllers

import (
	"encoding/json"
	"errors"
	"gym-api/src/auth"
	"gym-api/src/database"
	"gym-api/src/models"
	"gym-api/src/repository"
	"gym-api/src/responses"
	"gym-api/src/security"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.New(responses.ErrMsgUnprocessableEntity)
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	if err = user.Prepare("register"); err != nil {
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

	repo := repository.NewUsersRepository(db)
	user.ID, err = repo.Save(user)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	user.Password = ""

	responses.SendJSON(w, http.StatusCreated, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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
		err = errors.New(responses.ErrMsgForbidden)
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

	repo := repository.NewUsersRepository(db)
	userInDatabase, err := repo.FindByID(userID)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if userInDatabase.ID != userID {
		err = errors.New(responses.ErrMsgForbidden)
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.New(responses.ErrMsgUnprocessableEntity)
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if user.Name == "" {
		user.Name = userInDatabase.Name
	}
	if user.Email == "" {
		user.Email = userInDatabase.Email
	}

	if err = user.Prepare("update"); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	if err = repo.Update(userID, user); err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
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
		err = errors.New(responses.ErrMsgForbidden)
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

	repo := repository.NewUsersRepository(db)
	if err = repo.Delete(userID); err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
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
		err = errors.New(responses.ErrMsgForbidden)
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.New(responses.ErrMsgUnprocessableEntity)
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var password models.ChangePassword
	if err = json.Unmarshal(bodyRequest, &password); err != nil {
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

	repo := repository.NewUsersRepository(db)
	passwordInDatabase, err := repo.FindPassword(userID)
	if err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.VerifyPassword(passwordInDatabase, password.Current); err != nil {
		err = errors.New("the current password is invalid")
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	hashPassword, err := security.GenerateHash(password.New)
	if err != nil {
		err = errors.New(responses.ErrMsgBadRequest)
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	if err = repo.UpdatePassword(userID, string(hashPassword)); err != nil {
		err = errors.New(responses.ErrMsgInternalServerError)
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}
