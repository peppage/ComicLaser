package main

import (
	mdl "comiclaser/model"
	"comiclaser/monitor"
	"comiclaser/setting"

	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

func init() {
	setting.Initialize()
	log.SetLevel(7)
}

func main() {
	mdl.SetupDb()

	go monitor.Watch(setting.ComicFolder)

	e := echo.New()
	e.Get("/", root)
	e.Run(":8000")

}

func root(c *echo.Context) error {
	return c.String(http.StatusOK, "st")
}
