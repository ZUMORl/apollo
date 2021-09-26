package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// srv.GET("/devices/:d_id/sensors/:s_id/values", getValues)
func getValues(c echo.Context) error {
	var req = c.Request()
	if req.Header["Content-Type"][0] != "application/json" {
		return c.JSON(http.StatusBadRequest,
			fmt.Sprintf("%s is not accepted content type",
				req.Header["Content-Type"][0]))
	}

	var start, err = strconv.Atoi(req.Header["Start-Value-Index"][0])
	var end, err_1 = strconv.Atoi(req.Header["End-Value-Index"][0])
	if err != nil || err_1 != nil {
		return c.JSON(http.StatusBadRequest,
			fmt.Sprintf("Wrong headers values : %s, %s",
				req.Header["Start-Value-Index"][0], req.Header["End-Value-Index"][0]))
	}

	values, err := sensor.GetValues(c.Param("s_id"), start, end)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ServerError{
				err.Error(),
				fmt.Sprintf("GET /device/%v/sensors/%v/values",
					c.Param("d_id"), c.Param("s_id")),
				"",
			})
	}

	return c.JSON(http.StatusOK, values)
}

// srv.DELETE("/devices/:d_id/sensors/:s_id/values", deleteValues)
func deleteValues(c echo.Context) error {

	return nil
}

// srv.POST("/devices/:d_id/sensors/:s_id/values", newValue)
func newValue(c echo.Context) error {

	return nil
}
