package server

import (
	"fmt"
	"net/http"

	"github.com/apollo/db"

	"github.com/labstack/echo/v4"
)

func readSensor(c echo.Context) error {
	var dvc_id = c.Param("d_id")
	var ret = c.QueryParam("id")
	if ret == "" {
		var sensors, err = sensor.ListByDevice(dvc_id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		for key, elem := range sensors {
			ret += fmt.Sprintf("%v : %v\n", key, elem)
		}
	} else {
		var dvc, err = sensor.Read(ret)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		ret = fmt.Sprintf("%v : %v\n", ret, dvc)
	}
	return c.String(http.StatusOK, ret)
}

func newSensor(c echo.Context) error {
	var key, err = sensor.Add(&db.Sensor{
		Type:  c.Param("s_id"),
		Model: "default",
	}, c.Param("d_id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, key)
}
