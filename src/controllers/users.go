package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gym-api/src/auth"
	"gym-api/src/database"
	"gym-api/src/models"
	"gym-api/src/repository"
	"gym-api/src/responses"
	"gym-api/src/security"
	"gym-api/src/usecases"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserController struct {
	usecase *usecases.UserUseCase
}

func NewUsersController(usecase *usecases.UserUseCase) *UserController {
	return &UserController{usecase: usecase}
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	if err = user.Prepare("register"); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	user.ID, err = uc.usecase.Save(user)
	if err != nil {
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
		fmt.Println("erro aqui")
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	userIDInToken, err := auth.ExtractUserID(r)
	if err != nil {
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	if userID != userIDInToken {
		err = errors.New("cannot update a user other than your own")
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	if err = repo.Update(userID, user); err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
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
		err = errors.New("you cannot delete a user other than your own")
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	if err = repo.Delete(userID); err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
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
		err = errors.New("cannot change another user's password")
		responses.SendError(w, http.StatusForbidden, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.SendError(w, http.StatusUnprocessableEntity, err)
		return
	}

	var password models.ChangePassword
	if err = json.Unmarshal(bodyRequest, &password); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Connect()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	passwordInDatabase, err := repo.FindPassword(userID)
	if err != nil {
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
		responses.SendError(w, http.StatusBadRequest, err)
		return
	}

	if err = repo.UpdatePassword(userID, string(hashPassword)); err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	responses.SendJSON(w, http.StatusNoContent, nil)
}
