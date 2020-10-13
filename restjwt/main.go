// from https://www.udemy.com/course/build-jwt-authenticated-restful-apis-with-golang/learn/lecture/12604244

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"database/sql"
	"github.com/lib/pq"
	"os"
)

var db *sql.DB

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
	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db)

	db, err = sql.Open("postgres", pgUrl)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

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
