package server

import (
	"net/http"

	"github.com/apollo/db"

	"github.com/labstack/echo/v4"
)

var (
	device db.Devices
	sensor db.Sensors
)

type ServerError struct {
	Err       string
	Method    string
	ExtraInfo string
}

func home(c echo.Context) error {
	return c.String(http.StatusOK, "hello, world")
}

func Serve() {
	var srv = echo.New()

	device = db.NewDevices(db.Db)
	sensor = db.NewSensors(db.Db)

	srv.GET("/", home)

	// Routes for DEVICES
	srv.GET("/devices/", readDevices)
	srv.POST("/devices/", newDevice)
	srv.GET("/devices/:d_id/", readDevice)
	srv.PUT("/devices/:d_id/", updateDevice)
	srv.DELETE("/devices/:d_id/", deleteDevice)

	// Routes for SENSORS
	srv.GET("/devices/:d_id/sensors/", readSensors)
	srv.POST("/devices/:d_id/sensors/", newSensor)
	srv.GET("/devices/:d_id/sensors/:s_id/", readSensor)
	srv.PUT("/devices/:d_id/sensors/:s_id/", updateSensor)
	srv.DELETE("/devices/:d_id/sensors/:s_id/", deleteSensor)

	// Routes for VALUES
	srv.GET("/devices/:d_id/sensors/:s_id/values", getValues)
	srv.DELETE("/devices/:d_id/sensors/:s_id/values", deleteValues)
	srv.POST("/devices/:d_id/sensors/:s_id/values", newValue)

	srv.Logger.Fatal(srv.Start(":1323"))
}
