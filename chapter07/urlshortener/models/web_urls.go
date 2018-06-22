package models

import (
	"database/sql"

	_ "github.com/lib/pq" // Importing postgres driver
)

// InitDB returns an SQL connetion
func InitDB() (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", "postgresql://dakaraj:R71VDl6m@localhost:26257/mydb?sslmode=require")
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(`CREATE
TABLE IF NOT EXISTS web_url (
	id SERIAL PRIMARY KEY,
	url TEXT NOT NULL
);`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	return db, nil

}
