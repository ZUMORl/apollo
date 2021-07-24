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

func home(c echo.Context) error {

	return c.String(http.StatusOK, "hello, world")
}

func Serve() {
	var srv = echo.New()

	device = db.NewDevices(db.Db)
	sensor = db.NewSensors(db.Db)

	srv.GET("/", home)
	srv.GET("/device/new/:id", newDevice)
	srv.GET("/device/:d_id/sensor/new/:s_id", newSensor)
	srv.GET("/device", readDevice)
	srv.GET("/device/:d_id/sensor", readSensor)

	srv.Logger.Fatal(srv.Start(":1323"))
}
