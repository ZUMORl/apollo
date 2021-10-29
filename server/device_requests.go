package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/apollo/db"

	"github.com/labstack/echo/v4"
)

// srv.GET("/devices", readDevices)
func readDevices(c echo.Context) error {
	var devices, err = device.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				"GET /devices",
				"",
			})
	}

	return c.JSON(http.StatusOK, devices)
}

// srv.GET("/devices/:d_id", readDevice)
func readDevice(c echo.Context) error {
	var id = c.Param("d_id")

	var dvc, err = device.Read(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				fmt.Sprintf("GET /devices/%v", id),
				"",
			})
	}

	return c.JSON(http.StatusOK, dvc)
}

// srv.POST("/devices", newDevice)
func newDevice(c echo.Context) error {
	var req = c.Request()
	if req.Header["Content-Type"][0] != "application/json" {
		return c.JSON(http.StatusBadRequest,
			fmt.Sprintf("%s is not accepted content type",
				req.Header["Content-Type"][0]))
	}

	var newDevice = db.Device{}
	bits, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				"POST /devices",
				"",
			})
	}

	if err := json.Unmarshal(bits, &newDevice); err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				"POST /devices",
				"Incorrect json data. Could not decrypt.",
			})
	}

	key, err := device.Add(&newDevice)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				"POST /devices",
				"",
			})
	}
	return c.JSON(http.StatusOK, key)
}

// srv.PUT("/devices/:d_id", updateDevice)
func updateDevice(c echo.Context) error {
	var req = c.Request()
	var id = c.Param("d_id")
	if req.Header["Content-Type"][0] != "application/json" {
		return c.JSON(http.StatusBadRequest,
			fmt.Sprintf("%s is not accepted content type",
				req.Header["Content-Type"][0]))
	}

	var updatedDevice = db.Device{}
	bits, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				fmt.Sprintf("PUT /devices/%v", id),
				"",
			})
	}

	if err := json.Unmarshal(bits, &updatedDevice); err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				fmt.Sprintf("PUT /devices/%v", id),
				"Incorrect json data. Could not decrypt.",
			})
	}

	if err = device.Update(id, &updatedDevice); err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				fmt.Sprintf("PUT /devices/%v", id),
				fmt.Sprintf("name: %v, model: %v", updatedDevice.Name, updatedDevice.Model),
			})
	}

	return c.NoContent(http.StatusOK)
}

// srv.DELETE("/devices/:d_id", deleteDevice)
func deleteDevice(c echo.Context) error {
	var id = c.Param("d_id")

	if err := device.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				fmt.Sprintf("DELETE /devices/%v", id),
				"",
			})
	}

	return c.NoContent(http.StatusOK)
}
