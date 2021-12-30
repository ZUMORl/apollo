package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/apollo/db"

	"github.com/labstack/echo/v4"
)

// srv.GET("/devices/:d_id/sensors", readSensors)
func readSensors(c echo.Context) error {
	var dvc_id = c.Param("d_id")

	var sensors, err = sensor.ListByDevice(dvc_id)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message:  "Could not read Sensor information",
			Internal: err}
	}

	return c.JSON(http.StatusOK, sensors)
}

// srv.GET("/devices/:d_id/sensors/:s_id", readSensor)
func readSensor(c echo.Context) error {
	var id = c.Param("s_id")

	var sns, err = sensor.Read(id)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message:  "Could not read Sensor information",
			Internal: err}
	}

	return c.JSON(http.StatusOK, sns)
}

// srv.POST("/devices/:d_id/sensors", newSensor)
func newSensor(c echo.Context) error {
	var req = c.Request()
	if req.Header["Content-Type"][0] != "application/json" {
		return &echo.HTTPError{Code: http.StatusBadRequest,
			Message: fmt.Sprintf("%s is not accepted content type",
				req.Header["Content-Type"][0])}
	}

	var newSensor = db.Sensor{}
	bits, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message:  "Could not read request body",
			Internal: err}
	}

	if err := json.Unmarshal(bits, &newSensor); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message:  "Incorrect json data. Could not decrypt.",
			Internal: err}
	}

	key, err := sensor.Add(&newSensor, c.Param("d_id"))
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message:  "Could not add sensor to database",
			Internal: err}
	}
	return c.JSON(http.StatusOK, key)
}

// srv.PUT("/devices/:d_id/sensors/:s_id", updateSensor)
func updateSensor(c echo.Context) error {
	var req = c.Request()
	if req.Header["Content-Type"][0] != "application/json" {
		return &echo.HTTPError{Code: http.StatusBadRequest,
			Message: fmt.Sprintf("%s is not accepted content type",
				req.Header["Content-Type"][0])}
	}

	var updatedSensor = db.Sensor{}
	bits, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message:  "Could not read request body",
			Internal: err}
	}

	if err := json.Unmarshal(bits, &updatedSensor); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message:  "Incorrect json data. Could not decrypt.",
			Internal: err}
	}

	var id = c.Param("s_id")
	if err = sensor.Update(id, &updatedSensor); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message:  "Could not update sensor in database.",
			Internal: err}
	}

	return c.NoContent(http.StatusOK)
}

// srv.DELETE("/devices/:d_id/sensors/:s_id, deleteSensor)
func deleteSensor(c echo.Context) error {
	var id = c.Param("s_id")

	if err := sensor.Delete(id); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message:  "Could not delete sensor in database.",
			Internal: err}
	}

	return c.NoContent(http.StatusOK)
}
