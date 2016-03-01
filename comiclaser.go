package main

import (
	mdl "comiclaser/model"
	"comiclaser/monitor"
	"comiclaser/setting"

	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

func init() {
	setting.Initialize()
	ll, err := log.ParseLevel(strconv.Itoa(setting.LogLevel))
	if err == nil {
		log.SetLevel(ll)
	}
}

func main() {
	mdl.SetupDb(setting.DatabaseName)

	go monitor.Watch(setting.ComicFolder)

	e := echo.New()
	e.Get("/", root)
	e.Get("/dbinfo", dbinfo)
	log.Info("Server started on port " + setting.HttpPort)
	e.Run(":" + setting.HttpPort)

}

func root(c *echo.Context) error {
	return c.String(http.StatusOK, "st")
}

func dbinfo(c *echo.Context) error {
	return c.JSON(http.StatusOK, mdl.GetDbInfo())
}
