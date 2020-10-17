// from https://www.udemy.com/course/build-jwt-authenticated-restful-apis-with-golang/learn/lecture/12604244

package main

import (
	"database/sql"
	"github.com/golfz/learn-golang/restjwt/controllers"
	"github.com/golfz/learn-golang/restjwt/driver"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
)

var db *sql.DB

func init() {
	gotenv.Load()
	log.Println("init()")
}

func main() {
	db = driver.ConnectDB()

	ctrl := controllers.Controller{db}

	router := mux.NewRouter()

	router.HandleFunc("/signup", ctrl.Signup).Methods("POST")
	router.HandleFunc("/login", ctrl.Login).Methods("POST")
	router.HandleFunc("/protected", ctrl.TokenVerifyMiddleWare(controllers.ProtectedEndpoint)).Methods("POST")

	log.Println("Listen on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
