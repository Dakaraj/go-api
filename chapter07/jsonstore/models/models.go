package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // Driver for PostsgreSQL
)

// User struct contains users data table structure
type User struct {
	gorm.Model
	Order []Order
	Data  string `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB" json:""`
}

// Order struct contains orders data table structure
type Order struct {
	gorm.Model
	User User
	Data string `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
}

// InitDB initializes a connection to database
func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open("postgres", "postgresql://dakaraj:R71VDl6m@localhost:26257/mydb?sslmode=require")
	if err != nil {
		return nil, err
	}
	// The below AutoMigrate is equivalent to this
	// if !db.HasTable("users") {
	// 	db.CreateTable(&User{})
	// }
	// if !db.HasTable("orders") {
	// 	db.CreateTable(&Order{})
	// }

	db.AutoMigrate(&User{}, &Order{})

	return db, nil
}
