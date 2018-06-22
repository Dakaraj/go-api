package main

import (
	"encoding/json"
	"github.com/dakaraj/go-api/chapter07/jsonstore/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"time"
)

// DBClient stores database session information. Initialized once
type DBClient struct {
	db *gorm.DB
}

// UserResponse is the response to be sent back to client
type UserResponse struct {
	User models.User `json:"user"`
	Data interface{} `json:"data"`
}

func (driver *DBClient) getUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	vars := mux.Vars(r)
	// Handle response details
	driver.db.First(&user, vars["id"])
	var userData interface{}
	// Unmarshal JSON string to interface
	json.Unmarshal([]byte(user.Data), &userData)
	var response = UserResponse{User: user, Data: userData}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// responseMap := map[string]interface{}{"url": ""}
	respJSON, _ := json.Marshal(response)
	w.Write(respJSON)
}

func (driver *DBClient) postUser(w http.ResponseWriter, r *http.Request) {

}

func (driver *DBClient) getUsersByFirstName(w http.ResponseWriter, r *http.Request) {

}

func main() {
	db, err := models.InitDB()
	if err != nil {
		log.Fatal("DB initiation failed with error: ", err.Error())
	}

	dbclient := &DBClient{db: db}
	defer db.Close()

	r := mux.NewRouter()
	// Attach path with handler
	r.HandleFunc("/v1/user/{a-zA-Z0-9}+", dbclient.getUser).Methods("GET")
	r.HandleFunc("/v1/user", dbclient.postUser).Methods("POST")
	r.HandleFunc("/v1/user", dbclient.getUsersByFirstName).Methods("GET")
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8008",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
