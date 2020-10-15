// from https://www.udemy.com/course/build-jwt-authenticated-restful-apis-with-golang/learn/lecture/12604244

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"golang.org/x/crypto/bcrypt"
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

func responseError(w http.ResponseWriter, status int, err Error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func responseSuccessNoBody(w http.ResponseWriter, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
}

func responseSuccessWithBody(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func signup(w http.ResponseWriter, r *http.Request) {
	log.Println("signup invoked")

	var user User
	var error Error

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	if user.Email == "" {
		error.Message = "Email is empty"
		responseError(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "Password is empty"
		responseError(w, http.StatusBadRequest, error)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Password:", user.Password)
	fmt.Println("Hashed:", hashed)

	user.Password = string(hashed)

	fmt.Println("Password after hashed:", user.Password)


	stmt := "insert into users(email, password) values($1, $2) RETURNING id;"

	err = db.QueryRow(stmt, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		error.Message = "Create user error"
		responseError(w, http.StatusInternalServerError, error)
		return
	}

	user.Password = ""

	responseSuccessWithBody(w, http.StatusOK, user)

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
