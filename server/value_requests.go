package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/apollo/db"
	"github.com/labstack/echo/v4"
)

func compoundErrors(err_list []error) error {
	var str = ""
	for _, err := range err_list {
		if err != nil {
			str += err.Error() + ", "
		}
	}
	return errors.New(str)
}

// srv.GET("/devices/:d_id/sensors/:s_id/values", getValues)
func getValues(c echo.Context) error {
	var start, err = strconv.Atoi(c.QueryParam("start"))
	var end, err_1 = strconv.Atoi(c.QueryParam("end"))
	if err != nil || err_1 != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest,
			Message:  "Wrong query parameters",
			Internal: compoundErrors([]error{err, err_1})}
	}

	values, err := sensor.GetValues(c.Param("s_id"), start, end)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message: "Error retrieving sensor values", Internal: err}
	}

	return c.JSON(http.StatusOK, values)
}

// srv.DELETE("/devices/:d_id/sensors/:s_id/values", deleteValues)
func deleteValues(c echo.Context) error {
	var start, err = strconv.Atoi(c.QueryParam("start"))
	var end, err_1 = strconv.Atoi(c.QueryParam("end"))
	if err != nil || err_1 != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest,
			Message:  "Wrong query parameters",
			Internal: compoundErrors([]error{err, err_1})}
	}
	if err := sensor.RemoveValue(c.Param("s_id"), start, end); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message: "Error deleting sensor values", Internal: err}
	}
	return c.NoContent(http.StatusOK)
}

// srv.POST("/devices/:d_id/sensors/:s_id/values", newValue)
func newValue(c echo.Context) error {
	var newValue = db.Value{}
	bits, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message: "Error reading body", Internal: err}
	}

	if err := json.Unmarshal(bits, &newValue); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message: "Error parsing body json", Internal: err}
	}

	if err := sensor.AddValue(c.Param("s_id"), &newValue); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message: "Error adding value", Internal: err}
	}
	return c.NoContent(http.StatusOK)
}
