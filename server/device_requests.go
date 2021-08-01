package server

import (
	"fmt"
	"net/http"

	"github.com/apollo/db"

	"github.com/labstack/echo/v4"
)

// srv.GET("/devices/", readDevices)
func readDevices(c echo.Context) error {
	var ret string
	var devices, err = device.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				"GET /devices/",
				"",
			})
	}

	for key, elem := range devices {
		ret += fmt.Sprintf("%v : %v\n", key, elem)
	}

	return c.String(http.StatusOK, ret)
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

	var ret = fmt.Sprintf("%v : %v\n", id, dvc)
	return c.String(http.StatusOK, ret)
}

// srv.POST("/devices/", newDevice)
func newDevice(c echo.Context) error {
	var key, err = device.Add(&db.Device{
		Name:  c.FormValue("name"),
		Model: c.FormValue("model"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				"POST /devices/",
				fmt.Sprintf("name: %v, model: %v",
					c.FormValue("name"), c.FormValue("model")),
			})
	}
	return c.String(http.StatusOK, key)
}

// srv.PUT("/devices/:d_id", updateDevice)
func updateDevice(c echo.Context) error {
	var id = c.Param("d_id")
	var err = device.Update(id, &db.Device{
		Name:  c.FormValue("name"),
		Model: c.FormValue("model"),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				fmt.Sprintf("PUT /devices/%v", id),
				fmt.Sprintf("name: %v, model: %v",
					c.FormValue("name"), c.FormValue("model")),
			})
	}

	return c.String(http.StatusOK,
		fmt.Sprintf("%v Updated successfuly", id))
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

	var ret = fmt.Sprintf("%v Deleted successfuly\n", id)
	return c.String(http.StatusOK, ret)
}
