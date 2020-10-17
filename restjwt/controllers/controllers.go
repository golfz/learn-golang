package controllers

import "database/sql"

type Controller struct {
	Db *sql.DB
}
