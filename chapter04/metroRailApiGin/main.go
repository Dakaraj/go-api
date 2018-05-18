package main

import (
	"database/sql"
	"log"
	"regexp"

	"github.com/dakaraj/go-api/chapter04/dbutils"
	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

// DB driver that is visible for the whole program
var DB *sql.DB

var timePattern, _ = regexp.Compile(`^([01][0-9]|2[0-3]):[0-5][0-9]$`)

// TrainResource is a model for holding train information
type TrainResource struct {
	ID              int    `json:"id"`
	DriverName      string `json:"driverName"`
	OperatingStatus bool   `json:"operatingStatus"`
}

// GET /v1/trains/{train-id}
func getTrain(c *gin.Context) {
	var t TrainResource
	id := c.Param("train-id")
	err := DB.QueryRow(`
SELECT id, driver_name, operating_status
FROM train
WHERE id = ?;
`, id).Scan(&t.ID, &t.DriverName, &t.OperatingStatus)
	if err != nil {
		c.JSON(500, gin.H{
			"error":        err.Error(),
			"errorMessage": "Train could not be found",
		})
	} else {
		c.JSON(200, gin.H{
			"result": t,
		})
	}
}

// POST /v1/trains
func createTrain(c *gin.Context) {
	var t TrainResource
	err := c.BindJSON(&t)
	if err != nil || t.DriverName == "" {
		c.JSON(500, gin.H{
			"erorMessage": "Request JSON is invalid!",
		})
		return
	}

	statement, _ := DB.Prepare(`
INSERT INTO train
(driver_name, operating_status)
VALUES
(?, ?);
`)
	result, err := statement.Exec(t.DriverName, t.OperatingStatus)
	if err != nil {
		c.JSON(500, gin.H{
			"error":        err.Error(),
			"errorMessage": "Error inserting data into database",
		})
	} else {
		lastRowID, _ := result.LastInsertId()
		t.ID = int(lastRowID)
		c.JSON(201, gin.H{
			"result": t,
		})
	}
}

// PUT /v1/trains/
func updateTrain(c *gin.Context) {
	var t TrainResource
	err := c.BindJSON(&t)
	if err != nil || t.DriverName == "" || t.ID == 0 {
		c.JSON(500, gin.H{
			"errorMessage": "Request JSON is invalid!",
		})
		return
	}

	log.Printf("Train ID: %v, Driver name: %v, Operating status: %v", t.ID, t.DriverName, t.OperatingStatus)
	statement, _ := DB.Prepare(`
UPDATE train
SET driver_name = ?,
operating_status = ?
WHERE id = ?;
`)
	_, err = statement.Exec(t.DriverName, t.OperatingStatus, t.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":        err.Error(),
			"errorMessage": "Error updating row in database",
		})
	} else {
		c.JSON(200, gin.H{
			"result": t,
		})
	}
}

// DELETE /v1/trains/{train-id}
func removeTrain(c *gin.Context) {
	id := c.Param("train-id")
	statement, _ := DB.Prepare(`
DELETE FROM train
WHERE id = ?;
`)
	result, err := statement.Exec(id)
	if err != nil {
		c.JSON(500, gin.H{
			"error":        err.Error(),
			"errorMessage": "Error happened while attempting to delete row from database!",
		})
	} else {
		if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
			c.JSON(404, gin.H{
				"errorMessage": "Invalid train-id given. Train could not be deleted",
			})
		} else {
			c.JSON(200, map[string]string{"result": "success"})
		}
	}
}

// StationResource is a model for holding station information
type StationResource struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	OpeningTime string `json:"openingTime"`
	ClosingTime string `json:"closingTime"`
}

// GET /v1/stations/{station-id}
func getStation(c *gin.Context) {
	var s StationResource
	id := c.Param("station-id")
	err := DB.QueryRow(`
SELECT id, name, opening_time, closing_time
FROM station
WHERE id = ?;
`, id).Scan(&s.ID, &s.Name, &s.OpeningTime, &s.ClosingTime)
	if err != nil {
		c.JSON(500, gin.H{
			"error":        err.Error(),
			"errorMessage": "Station could not be found",
		})
	} else {
		c.JSON(200, gin.H{
			"result": s,
		})
	}
}

// POST /v1/stations
func createStation(c *gin.Context) {
	var s StationResource
	err := c.BindJSON(&s)
	if err != nil {
		c.JSON(500, gin.H{
			"error":       err.Error(),
			"erorMessage": "Request JSON is invalid!",
		})
		return
	} else if s.Name == "" {
		c.JSON(500, gin.H{
			"erorMessage": "Please provide a station name!",
		})
		return
	} else if !timePattern.MatchString(s.OpeningTime) || !timePattern.MatchString(s.ClosingTime) {
		c.JSON(500, gin.H{
			"errorMessage": `Provided time format is invalid. Use "HH:mm"!`,
		})
		return
	}
	statement, _ := DB.Prepare(`
INSERT INTO station
(name, opening_time, closing_time)
VALUES (?, ?, ?);
`)
	result, err := statement.Exec(s.Name, s.OpeningTime, s.ClosingTime)
	if err != nil {
		c.JSON(500, gin.H{
			"error":        err.Error(),
			"errorMessage": "Error inserting data into database",
		})
	} else {
		lastRowID, _ := result.LastInsertId()
		s.ID = int(lastRowID)
		c.JSON(201, gin.H{
			"result": s,
		})
	}
}

// PUT /v1/stations
func updateStation(c *gin.Context) {
	var s StationResource
	err := c.BindJSON(&s)
	if err != nil {
		c.JSON(500, gin.H{
			"error":       err.Error(),
			"erorMessage": "Request JSON is invalid!",
		})
		return
	} else if s.Name == "" || s.ID == 0 {
		c.JSON(500, gin.H{
			"erorMessage": "Please provide a station name and Id!",
		})
		return
	} else if !timePattern.MatchString(s.OpeningTime) || !timePattern.MatchString(s.ClosingTime) {
		c.JSON(500, gin.H{
			"errorMessage": `Provided time format is invalid. Use "HH:mm"!`,
		})
		return
	}
	statement, _ := DB.Prepare(`
UPDATE station
SET name = ?,
opening_time = ?,
closing_time = ?
WHERE id = ?;
`)
	_, err = statement.Exec(s.Name, s.OpeningTime, s.ClosingTime, s.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"error":        err.Error(),
			"errorMessage": "Error inserting data into database",
		})
	} else {
		c.JSON(201, gin.H{
			"result": s,
		})
	}
}

// DELETE /v1/stations/{station-id}
func removeStation(c *gin.Context) {
	id := c.Param("station-id")
	statement, _ := DB.Prepare(`
DELETE FROM station
WHERE id = ?;
`)
	result, err := statement.Exec(id)
	if err != nil {
		c.JSON(500, gin.H{
			"error":        err.Error(),
			"errorMessage": "Error happened while attempting to delete row from database!",
		})
	} else {
		if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
			c.JSON(404, gin.H{
				"errorMessage": "Invalid station-id given. Station could not be deleted",
			})
		} else {
			c.JSON(200, map[string]string{"result": "success"})
		}
	}
}

// ScheduleResource is a model for holding station information
type ScheduleResource struct {
	ID          int
	TrainID     int
	StationID   int
	ArrivalTime string
}

func main() {
	var err error
	DB, err = sql.Open("sqlite3", "./chapter04/railapi.db")
	if err != nil {
		log.Fatal("Table creation failed:", err)
	}
	dbutils.Initialize(DB)

	router := gin.Default()
	router.GET("/v1/trains/:train-id", getTrain)
	router.POST("/v1/trains", createTrain)
	router.PUT("/v1/trains", updateTrain)
	router.DELETE("/v1/trains/:train-id", removeTrain)

	router.GET("/v1/stations/:station-id", getStation)
	router.POST("/v1/stations", createStation)
	router.PUT("/v1/stations", updateStation)
	router.DELETE("/v1/stations/:train-id", removeStation)

	router.Run(":8080")
}
