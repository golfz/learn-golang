package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/golfz/learn-golang/restjwt/models"
	userRepository "github.com/golfz/learn-golang/restjwt/repository/user"
	"github.com/golfz/learn-golang/restjwt/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strings"
)

func (ctrl Controller) Signup(w http.ResponseWriter, r *http.Request) {
	log.Println("signup invoked")

	var user models.User
	var error models.Error

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	if user.Email == "" {
		error.Message = "Email is empty"
		utils.ResponseError(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "Password is empty"
		utils.ResponseError(w, http.StatusBadRequest, error)
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

	userRepo := userRepository.UserRepository{}
	user = userRepo.Signup(ctrl.Db, user)

	utils.ResponseSuccessWithBody(w, http.StatusOK, user)

}

func (ctrl Controller) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("login invoked")

	var user models.User
	var error models.Error
	var jwt models.JWT

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	if user.Email == "" {
		error.Message = "Email is empty"
		utils.ResponseError(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "Password is empty"
		utils.ResponseError(w, http.StatusBadRequest, error)
		return
	}

	password := user.Password

	log.Println("password:", password)

	userRepo := userRepository.UserRepository{}
	user, err = userRepo.Login(ctrl.Db, user)

	if err != nil {
		error.Message = "Cannot login"
		utils.ResponseError(w, http.StatusUnauthorized, error)
		return
	}

	spew.Dump(user)

	hashedPassword := user.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		error.Message = "password is not corrected"
		utils.ResponseError(w, http.StatusUnauthorized, error)
		return
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		error.Message = "cannot create a token"
		utils.ResponseError(w, http.StatusInternalServerError, error)
		return
	}

	log.Println(token)

	jwt.Token = token

	utils.ResponseSuccessWithBody(w, http.StatusOK, jwt)
}

func (ctrl Controller) TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	log.Println("TokenVerifyMiddleWare invoked")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObj models.Error

		authHeader := r.Header.Get("Authorization")

		log.Println(authHeader)

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			errorObj.Message = "Wrong token"
			utils.ResponseError(w, http.StatusForbidden, errorObj)
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
			utils.ResponseError(w, http.StatusUnauthorized, errorObj)
			return
		}

		spew.Dump(token)

		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			errorObj.Message = "token is invalid"
			utils.ResponseError(w, http.StatusUnauthorized, errorObj)
			return
		}

	})
}
