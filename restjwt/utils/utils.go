package utils

import (
	"encoding/json"
	"github.com/golfz/learn-golang/restjwt/models"
	"net/http"
)

func ResponseError(w http.ResponseWriter, status int, err models.Error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func ResponseSuccessWithBody(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
