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
	e.Get("/comic/:id", getComic)
	e.Get("/comic/:id/page/:page", getPage)
	log.Info("Server started on port " + setting.HttpPort)
	e.Run(":" + setting.HttpPort)

}

func root(c *echo.Context) error {
	return c.String(http.StatusOK, "st")
}

func dbinfo(c *echo.Context) error {
	return c.JSON(http.StatusOK, mdl.GetDbInfo())
}

func getComic(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    c.Param("id"),
		}).Error("Invalid URL ID")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}
	comics, err := mdl.GetComic(id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    c.Param("id"),
		}).Error("Unable to get comic")
	}
	data := struct {
		Comics *[]mdl.Comic `json:"comics"`
		Total  int          `json:"total_count"`
		Pages  int          `json:"page_count"`
	}{
		comics,
		1,
		1,
	}
	return c.JSON(http.StatusOK, data)
}

func getPage(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    c.Param("id"),
		}).Error("Invalid URL ID")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"page":  c.Param("page"),
		}).Error("Invalid URL page")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Page")
	}

	p, err := mdl.GetComicPage(id, page)

	if err != nil {
		log.WithError(err).Error("Problem getting page from db or archive")
		return echo.NewHTTPError(http.StatusBadRequest, "Page doesn't exist")
	}

	c.Response().Header().Set("Content-Type", "image/jpeg")
	c.Response().Write(p)
	return nil
}
