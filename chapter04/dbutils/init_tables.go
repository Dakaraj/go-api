package dbutils

import (
	"database/sql"
	"log"
)

// Initialize function creates tables in DB based on schemas provided in schemas.go
func Initialize(dbDriver *sql.DB) {
	statement, driverError := dbDriver.Prepare(trainTable)
	if driverError != nil {
		log.Println(driverError)
	}
	statement.Exec()
	statement, _ = dbDriver.Prepare(stationTable)
	statement.Exec()
	statement, _ = dbDriver.Prepare(scheduleTable)
	statement.Exec()
	log.Println("All tables created/initialized successfully!")
}
