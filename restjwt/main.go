// from https://www.udemy.com/course/build-jwt-authenticated-restful-apis-with-golang/learn/lecture/12604244

package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/signup", signup).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/protected", TokenVerifyMiddleWare(ProtectedEndpoint)).Methods("POST")

	log.Println("Listen on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func signup(w http.ResponseWriter, r *http.Request) {
	log.Println("signup invoked")
	w.Write([]byte("hello this is Signup"))
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("login invoked")
}

func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("ProtectedEndpoint invoked")
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	log.Println("TokenVerifyMiddleWare invoked")
	return nil
}
