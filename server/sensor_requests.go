package server

import (
	"fmt"
	"net/http"

	"github.com/apollo/db"

	"github.com/labstack/echo/v4"
)

// srv.GET("/device/:d_id/sensors/", readSensors)
func readSensors(c echo.Context) error {
	var dvc_id = c.Param("d_id")
	var ret string

	var sensors, err = sensor.ListByDevice(dvc_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				fmt.Sprintf("GET /device/%v/sensors/", c.Param("d_id")),
				"",
			})
	}

	for key, elem := range sensors {
		ret += fmt.Sprintf("%v : %v\n", key, elem)
	}
	return c.String(http.StatusOK, ret)
}

// srv.GET("/device/:d_id/sensors/:s_id", readSensor)
func readSensor(c echo.Context) error {
	var id = c.Param("s_id")

	var sns, err = sensor.Read(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				fmt.Sprintf("GET /device/%v/sensors/%v",
					c.Param("d_id"), c.Param("s_id")),
				"",
			})
	}
	var ret = fmt.Sprintf("%v : %v\n", id, sns)

	return c.String(http.StatusOK, ret)
}

// srv.POST("/device/:d_id/sensors/", newSensor)
func newSensor(c echo.Context) error {
	var key, err = sensor.Add(&db.Sensor{
		Type:  c.FormValue("type"),
		Model: c.FormValue("model"),
	}, c.Param("d_id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				fmt.Sprintf("POST /device/%v/sensors/",
					c.Param("d_id")),
				fmt.Sprintf("type: %v, model: %v",
					c.FormValue("type"), c.FormValue("model")),
			})
	}

	return c.String(http.StatusOK, key)
}

// srv.PUT("/device/:d_id/sensors/:s_id", updateSensor)
func updateSensor(c echo.Context) error {
	var id = c.Param("s_id")
	var err = sensor.Update(id, &db.Sensor{
		Type:  c.FormValue("type"),
		Model: c.FormValue("model"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				fmt.Sprintf("PUT /device/%v/sensors/%v",
					c.Param("d_id"), c.Param("s_id")),
				fmt.Sprintf("type: %v, model: %v",
					c.FormValue("type"), c.FormValue("model")),
			})
	}

	return c.String(http.StatusOK,
		fmt.Sprintf("%v Updated successfuly", id))
}

// srv.DELETE("/device/:d_id/sensors/:s_id", deleteSensor)
func deleteSensor(c echo.Context) error {
	var id = c.Param("s_id")

	if err := sensor.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				fmt.Sprintf("DELETE /device/%v/sensors/%v",
					c.Param("d_id"), c.Param("s_id")),
				"",
			})
	}

	var ret = fmt.Sprintf("%v Deleted successfuly\n", id)
	return c.String(http.StatusOK, ret)
}
