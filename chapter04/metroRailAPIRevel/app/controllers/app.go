package controllers

import (
	"strconv"

	"github.com/revel/revel"
)

// App devines a general struct for application
type App struct {
	*revel.Controller
}

// TrainResource is the model for holding rail information
type TrainResource struct {
	ID              int    `json:"id"`
	DriverName      string `json:"driverName"`
	OperatingStatus bool   `json:"operatingStatus"`
}

// GetTrain handles GET on train resource: /v1/trains/:train-id
func (c App) GetTrain() revel.Result {
	var t TrainResource
	id := c.Params.Route.Get("train-id")
	// 	err := DB.QueryRow(`
	// SELECT id, driver_name, operating_status
	// FROM train
	// WHERE id = ?;
	// `, id).Scan(&t.ID, &t.DriverName, &t.OperatingStatus)

	// HERE COMES DB MOCKING
	intID, err := strconv.Atoi(id)
	t.ID = intID
	t.DriverName = "Daka"
	t.OperatingStatus = true
	// END OF DB MOCKING

	if err != nil {
		return c.RenderJSON(map[string]string{
			"error":        err.Error(),
			"errorMessage": "Train could not be found",
		})
	} else {
		return c.RenderJSON(t)
	}
}
