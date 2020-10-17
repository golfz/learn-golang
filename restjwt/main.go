// from https://www.udemy.com/course/build-jwt-authenticated-restful-apis-with-golang/learn/lecture/12604244

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/subosito/gotenv"
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

func init() {
	gotenv.Load()
	log.Println("init()")
}

func main() {
	log.Println("db connection:", os.Getenv("ELEPHANTSQL_URL"))
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

func GenerateToken(user User) (string, error) {
	var err error
	secret := os.Getenv("SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "course",
	})

	spew.Dump(token)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("login invoked")

	var user User
	var error Error
	var jwt JWT

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

	password := user.Password

	log.Println("password:", password)

	stmt := "select * from users where email=$1"
	row := db.QueryRow(stmt, user.Email)

	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		log.Println("error", err)
	}

	spew.Dump(user)

	hashedPassword := user.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		error.Message = "password is not corrected"
		responseError(w, http.StatusUnauthorized, error)
		return
	}

	token, err := GenerateToken(user)
	if err != nil {
		error.Message = "cannot create a token"
		responseError(w, http.StatusInternalServerError, error)
		return
	}

	log.Println(token)

	jwt.Token = token

	responseSuccessWithBody(w, http.StatusOK, jwt)
}

func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("ProtectedEndpoint invoked")
	responseSuccessWithBody(w, http.StatusOK, "this is protected data")
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	log.Println("TokenVerifyMiddleWare invoked")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObj Error

		authHeader := r.Header.Get("Authorization")

		log.Println(authHeader)

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			errorObj.Message = "Wrong token"
			responseError(w, http.StatusForbidden, errorObj)
			return
		}

		authToken := bearerToken[1]
		log.Println(authToken)

		token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}

			return []byte(os.Getenv("SECRET")), nil
		})
		if error != nil {
			errorObj.Message = error.Error()
			responseError(w, http.StatusUnauthorized, errorObj)
			return
		}

		spew.Dump(token)

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			errorObj.Message = "token is invalid"
			responseError(w, http.StatusUnauthorized, errorObj)
			return
		}

	})
}
