package controllers

import (
	"encoding/json"
	"gym-api/src/auth"
	"gym-api/src/database"
	"gym-api/src/models"
	"gym-api/src/repository"
	"gym-api/src/responses"
	"gym-api/src/security"
	"io"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
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

	db, err := database.Connect()
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewUsersRepository(db)
	userInDatabase, err := repo.GetByEmail(user.Email)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.VerifyPassword(userInDatabase.Password, user.Password); err != nil {
		responses.SendError(w, http.StatusUnauthorized, err)
		return
	}

	token, err := auth.GenerateToken(userInDatabase.ID)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err)
		return
	}

	w.Write([]byte(token))
}
