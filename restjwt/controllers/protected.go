package controllers

import (
	"github.com/golfz/learn-golang/restjwt/utils"
	"log"
	"net/http"
)

func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("ProtectedEndpoint invoked")
	utils.ResponseSuccessWithBody(w, http.StatusOK, "this is protected data")
}
