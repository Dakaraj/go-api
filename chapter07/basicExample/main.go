package main

import (
	"github.com/dakaraj/go-api/chapter07/basicExample/models"
	"log"
)

func main() {
	db, err := models.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(db)
}
