package middlewares

import (
	"gym-api/src/auth"
	"gym-api/src/responses"
	"net/http"
)

func VerifyAuthenticateUser(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := auth.ValidateToken(r); err != nil {
			responses.SendError(w, http.StatusUnauthorized, err)
			return
		}
		nextFunction(w, r)
	}
}
