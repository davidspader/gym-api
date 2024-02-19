package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Fatal(err)
		}
	}
}

func SendError(w http.ResponseWriter, statusCode int, err error) {
	SendJSON(w, statusCode, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}
