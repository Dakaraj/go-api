package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Book struct used for retreiving book data from DB
type Book struct {
	id     int
	name   string
	author string
}

func main() {
	db, err := sql.Open("sqlite3", "./chapter04/sqliteFundamentals/books.db")
	log.Println(db)
	if err != nil {
		log.Println(err)
	}

	// Dropping table
	statement, err := db.Prepare(`
DROP TABLE
IF EXISTS books;
`)
	if err != nil {
		log.Println("Error while dropping table:", err)
	} else {
		log.Println("Table 'books' cleared!")
	}

	// Create table
	statement, _ = db.Prepare(`
CREATE TABLE
IF NOT EXISTS books (
	id INTEGER PRIMARY KEY,
	isbn INTEGER,
	author VARCHAR(64),
	name VARCHAR(64) NULL
);`)
	if err != nil {
		log.Println("Error encountered while creating table!")
	} else {
		log.Println("Table 'books' successfully created!")
	}
	statement.Exec()

	// Create
	statement, _ = db.Prepare(`
INSERT INTO books (
	name, author, isbn
)
VALUES (?, ?, ?);
`)
	statement.Exec("A Tale Of Two Cities", "Charles Dickens", 140430547)
	log.Println("Inserted the book into DB!")

	// Read
	rows, _ := db.Query(`
SELECT
id, name, author
FROM books;
`)
	var tempBook Book
	for rows.Next() {
		rows.Scan(&tempBook.id, &tempBook.name, &tempBook.author)
		log.Printf("ID: %d; Name: %s; Author: %s\n", tempBook.id, tempBook.name, tempBook.author)
	}

	// Update
	statement, _ = db.Prepare(`
UPDATE books
SET name = ?
WHERE id = ?;
`)
	statement.Exec("The Tale Of Two Cities", 1)
	log.Println("Successfuly updated the book in DB!")

	// Delete
	statement, _ = db.Prepare(`
DELETE FROM books
WHERE id = ?;
`)
	statement.Exec(1)
	log.Println("Successfuly deleted the book from DB!")
}
