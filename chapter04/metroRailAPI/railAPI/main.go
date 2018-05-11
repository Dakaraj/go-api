package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dakaraj/go-api/chapter04/metroRailAPI/dbutils"
	"github.com/emicklei/go-restful"
	_ "github.com/mattn/go-sqlite3"
)

// DB driver that is visible for the whole program
var DB *sql.DB

// TrainResource is a model for holding train information
type TrainResource struct {
	ID              int
	DriverName      string
	OperatingStatus bool
}

// Register adds paths and routes to container
func (t *TrainResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/trains").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{train-id}").To(t.getTrain))
	ws.Route(ws.POST("").To(t.createTrain))
	ws.Route(ws.PUT("").To(t.updateTrain))
	ws.Route(ws.DELETE("/{train-id}").To(t.removeTrain))
	container.Add(ws)
}

// GET /v1/trains/{train-id}
func (t *TrainResource) getTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train-id")
	err := DB.QueryRow(`
SELECT id, driver_name, operating_status
FROM train
WHERE id = ?;
`, id).Scan(&t.ID, &t.DriverName, &t.OperatingStatus)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Train could not be found.")
	} else {
		response.WriteEntity(t)
	}
}

// POST /v1/trains
func (t *TrainResource) createTrain(request *restful.Request, response *restful.Response) {
	log.Println(request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var b TrainResource
	err := decoder.Decode(&b)
	log.Printf("Driver name: %v, Operating status: %v", b.DriverName, b.OperatingStatus)
	statement, _ := DB.Prepare(`
INSERT INTO train
(driver_name, operating_status)
VALUES
(?, ?);
`)
	result, err := statement.Exec(b.DriverName, b.OperatingStatus)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	} else {
		lastRowID, _ := result.LastInsertId()
		b.ID = int(lastRowID)
		response.WriteHeaderAndEntity(http.StatusCreated, b)
	}
}

// PUT /v1/trains/
func (t *TrainResource) updateTrain(request *restful.Request, response *restful.Response) {
	log.Println(request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var b TrainResource
	err := decoder.Decode(&b)
	log.Printf("Train ID: %v, Driver name: %v, Operating status: %v", b.ID, b.DriverName, b.OperatingStatus)
	statement, _ := DB.Prepare(`
UPDATE train
SET driver_name = ?,
operating_status = ?
WHERE id = ?;
`)
	_, err = statement.Exec(b.DriverName, b.OperatingStatus, b.ID)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, b)
	}
}

// DELETE /v1/trains/{train-id}
func (t *TrainResource) removeTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train-id")
	statement, _ := DB.Prepare(`
DELETE FROM train
WHERE id = ?;
`)
	_, err := statement.Exec(id)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteHeader(http.StatusOK)
		response.Write([]byte("Train deleted successfully"))
	}
}

// StationTime is a wrapper for time.Time to override several methods
type StationTime struct {
	time.Time
}

// Scan used for parsing 'time' from database
func (st *StationTime) Scan(i interface{}) error {
	bs, _ := i.([]byte)
	val := string(bs)
	parsedTime, err := time.Parse("15:04", val)
	if err != nil {
		return err
	}
	st.Time = parsedTime
	return nil
}

// MarshalJSON is used to marshal time.Time format to HH:MM format
func (st StationTime) MarshalJSON() ([]byte, error) {
	stringTime := fmt.Sprintf("\"%s\"", st.Format("15:04"))
	return []byte(stringTime), nil
}

// UnmarshalJSON is used to parse HH:MM formatted string to time.Time format
func (st *StationTime) UnmarshalJSON(data []byte) error {
	stringTime := string(data)
	parsedTime, err := time.Parse("\"15:04\"", stringTime)
	if err != nil {
		return err
	}
	st.Time = parsedTime
	return nil
}

// StationResource is a model for holding station information
type StationResource struct {
	ID          int
	Name        string
	OpeningTime StationTime
	ClosingTime StationTime
}

// Register adds paths and routes to container
func (s *StationResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/stations").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{station-id}").To(s.getStation))
	ws.Route(ws.POST("").To(s.createStation))
	ws.Route(ws.PUT("").To(s.updateStation))
	ws.Route(ws.DELETE("/{station-id}").To(s.removeStation))
	container.Add(ws)
}

// GET /v1/stations/{station-id}
func (s *StationResource) getStation(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("station-id")
	log.Printf("Station id: %v\n", id)
	err := DB.QueryRow(`
SELECT id, name, opening_time, closing_time
FROM station
WHERE id = ?;
`, id).Scan(&s.ID, &s.Name, &s.OpeningTime, &s.ClosingTime)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Station could not be found.")
	} else {
		response.WriteEntity(s)
	}
}

// POST /v1/station
func (s *StationResource) createStation(request *restful.Request, response *restful.Response) {
	log.Println(request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var b StationResource
	err := decoder.Decode(&b)
	log.Printf("Station name: %v, opening time: %v, closing time: %v\n",
		b.Name, b.OpeningTime, b.ClosingTime)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	statement, _ := DB.Prepare(`
INSERT INTO station
(name, opening_time, closing_time)
VALUES (?, ?, ?);
`)
	result, err := statement.Exec(b.Name, b.OpeningTime.Format("15:04"), b.ClosingTime.Format("15:04"))
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())

	} else {
		lastRowID, _ := result.LastInsertId()
		b.ID = int(lastRowID)
		response.WriteHeaderAndEntity(http.StatusCreated, b)
	}
}

// PUT /v1/stations
func (s *StationResource) updateStation(request *restful.Request, response *restful.Response) {
	log.Println(request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var b StationResource
	err := decoder.Decode(&b)
	log.Printf("Station ID: %v, Station name: %v, Opening time: %v, Closing time: %v",
		b.ID, b.Name, b.OpeningTime, b.ClosingTime)
	statement, _ := DB.Prepare(`
UPDATE station
SET name = ?,
opening_time = ?,
closing_time = ?
WHERE id = ?;
`)
	_, err = statement.Exec(b.Name, b.OpeningTime.Format("15:04"), b.ClosingTime.Format("15:04"), b.ID)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, b)
	}
}

// DELETE /v1/stations/{station-id}
func (s *StationResource) removeStation(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("station-id")
	statement, _ := DB.Prepare(`
DELETE FROM station
WHERE id = ?;
`)
	_, err := statement.Exec(id)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteHeader(http.StatusOK)
		response.Write([]byte("Station deleted successfully"))
	}
}

// ScheduleResource is a model for holding station information
type ScheduleResource struct {
	ID          int
	TrainID     int
	StationID   int
	ArrivalTime StationTime
}

// Register adds paths and routes to container
func (s *ScheduleResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/schedules").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{schedule-id}").To(s.getSchedule))
	ws.Route(ws.POST("").To(s.createSchedule))
	ws.Route(ws.PUT("").To(s.updateSchedule))
	ws.Route(ws.DELETE("/{schedule-id}").To(s.removeSchedule))
	container.Add(ws)
}

// GET /v1/schedules/{schedule-id}
func (s *ScheduleResource) getSchedule(request *restful.Request, response *restful.Response) {
	scheduleID := request.PathParameter("schedule-id")
	log.Println("Schedule id:", scheduleID)
	err := DB.QueryRow(`
SELECT id, train_id, station_id, arrival_time
FROM schedule
WHERE id = ?;
`, scheduleID).Scan(&s.ID, &s.TrainID, &s.StationID, &s.ArrivalTime)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Schedule can not be found.")
	} else {
		response.WriteEntity(s)
	}
}

// POST /v1/schedules
func (s *ScheduleResource) createSchedule(request *restful.Request, response *restful.Response) {
	log.Println(request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var b ScheduleResource
	err := decoder.Decode(&b)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Train ID: %v, Station id: %v, Arrival time: %v\n",
		b.TrainID, b.StationID, b.ArrivalTime.Format("15:04"))
	statement, _ := DB.Prepare(`
INSERT INTO schedule
(train_id, station_id, arrival_time)
VALUES (?, ?, ?);
`)
	result, err := statement.Exec(b.TrainID, b.StationID, b.ArrivalTime.Format("15:04"))
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	} else {
		lastRowID, _ := result.LastInsertId()
		b.ID = int(lastRowID)
		response.WriteHeaderAndEntity(http.StatusCreated, b)
	}
}

// PUT /v1/schedules
func (s *ScheduleResource) updateSchedule(request *restful.Request, response *restful.Response) {
	log.Println(request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var b ScheduleResource
	err := decoder.Decode(&b)
	log.Printf("Schedule ID: %v, Train ID: %v, Station ID: %v, Arrival time: %v",
		b.ID, b.TrainID, b.StationID, b.ArrivalTime)
	statement, _ := DB.Prepare(`
UPDATE schedule
SET train_id = ?,
station_id = ?,
arrival_time = ?
WHERE id = ?;
`)
	_, err = statement.Exec(b.TrainID, b.StationID, b.ArrivalTime.Format("15:04"), b.ID)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, b)
	}
}

// DELETE /v1/schedules/{schedule-id}
func (s *ScheduleResource) removeSchedule(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("schedule-id")
	log.Printf("Schedule id: %v\n", id)
	statement, _ := DB.Prepare(`
DELETE FROM schedule
WHERE id = ?;
`)
	_, err := statement.Exec(id)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteHeader(http.StatusOK)
		response.Write([]byte("Schedule deleted successfully."))
	}
}

func main() {
	var err error
	DB, err = sql.Open("sqlite3", "./chapter04/metroRailAPI/railapi.db")
	if err != nil {
		log.Fatal("Table creation failed:", err)
	}
	dbutils.Initialize(DB)
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	// Register TrainResource
	t := TrainResource{}
	t.Register(wsContainer)
	// Register StationResource
	s := StationResource{}
	s.Register(wsContainer)
	// Register ScheduleResource
	sr := ScheduleResource{}
	sr.Register(wsContainer)
	log.Println("Start listening on localhost:80")
	server := &http.Server{
		Addr:    ":80",
		Handler: wsContainer,
	}
	log.Fatal(server.ListenAndServe())
}
