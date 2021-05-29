package service

import (
	"encoding/json"
	"log"
	"net/http"
)

func JsonError(w http.ResponseWriter, msg string, statusCode int) {
	var err struct {
		ErrorMessage string `json:"error_message"`
		StatusCode   int    `json:"status_code"`
	}

	err.ErrorMessage = msg
	err.StatusCode = statusCode

	if err := json.NewEncoder(w).Encode(&err); err != nil {
		log.Println("error sending error message...")
	}
}
