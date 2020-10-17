package userRepository

import (
	"database/sql"
	"github.com/golfz/learn-golang/restjwt/models"
	"log"
)

type UserRepository struct {

}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (u UserRepository) Signup(db *sql.DB, user models.User) models.User {
	stmt := "insert into users(email, password) values($1, $2) RETURNING id;"

	err := db.QueryRow(stmt, user.Email, user.Password).Scan(&user.ID)

	logFatal(err)

	user.Password = ""

	return user
}

func (u UserRepository) Login(db *sql.DB, user models.User) (models.User, error) {
	stmt := "select * from users where email=$1"
	row := db.QueryRow(stmt, user.Email)

	err := row.Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
