package main

import (
	"encoding/json"
	"log"

	"github.com/dakaraj/go-api/chapter07/jsonstore/models"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

// DBClient stores database session information. Initialized once
type DBClient struct {
	db *gorm.DB
}

// UserResponse is the response to be sent back to client
type UserResponse struct {
	Success  bool     `json:"success"`
	ID       uint     `json:"id"`
	UserData UserData `json:"userData"`
}

// UserData is a JSON struct that represent a user personal data
type UserData struct {
	EmailAddress string `json:"email_address"`
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

func (driver *DBClient) getUser(ctx iris.Context) {
	var user = models.User{}
	userID := ctx.Params().GetDecoded("id")
	// Handle response details
	driver.db.First(&user, userID)
	if user.ID == 0 {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"success": false, "message": "User not found"})
	} else {
		var userData UserData
		// Unmarshal JSON string to interface
		json.Unmarshal([]byte(user.Data), &userData)
		ctx.JSON(UserResponse{Success: true, ID: user.ID, UserData: userData})
	}
}

func (driver *DBClient) postUser(ctx iris.Context) {
	var user models.User
	var body UserData
	err := ctx.ReadJSON(&body)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"success": false, "message": "Invalid JSON in the request"})
	} else {
		bytesBody, _ := json.Marshal(body)
		user.Data = string(bytesBody)
		driver.db.Save(&user)
		ctx.JSON(iris.Map{"success": true, "userId": user.ID})
	}
	
}

func (driver *DBClient) getUsersByFirstName(ctx iris.Context) {

}

func main() {
	db, err := models.InitDB()
	if err != nil {
		log.Fatal("DB initiation failed with error: ", err.Error())
	}

	dbclient := &DBClient{db: db}
	defer db.Close()

	app := iris.New()
	app.Logger().SetLevel("debug")

	app.Use(recover.New())
	app.Use(logger.New())

	// Attach path with handler
	app.Get(`/v1/user/{id:string regexp(\d+)}`, dbclient.getUser)
	app.Post("/v1/user", dbclient.postUser)
	app.Get("/v1/user", dbclient.getUsersByFirstName)

	app.Run(iris.Addr(":8008"), iris.WithoutServerError(iris.ErrServerClosed))
}
