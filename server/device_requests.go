package server

import (
	"fmt"
	"net/http"

	"github.com/apollo/db"

	"github.com/labstack/echo/v4"
)

func readDevice(c echo.Context) error {
	var ret = c.QueryParam("id")
	if ret == "" {
		var devices, err = device.List()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		for key, elem := range devices {
			ret += fmt.Sprintf("%v : %v\n", key, elem)
		}
	} else {
		var dvc, err = device.Read(ret)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		ret = fmt.Sprintf("%v : %v\n", ret, dvc)
	}
	return c.String(http.StatusOK, ret)
}

func newDevice(c echo.Context) error {
	var key, err = device.Add(&db.Device{
		Name:  c.Param("id"),
		Model: "default",
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, key)
}
