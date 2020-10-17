package driver

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func ConnectDB() *sql.DB {
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

	return db
}
