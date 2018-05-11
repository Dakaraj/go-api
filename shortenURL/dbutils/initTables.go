package dbutils

import (
	"database/sql"
	"log"
)

// Initialize all SQL tables
func Initialize(dbDriver *sql.DB) {
	statement, err := dbDriver.Prepare(shortenURLTable)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	log.Println("Tables initiated successfully.")
}
