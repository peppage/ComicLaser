package main

import (
	mdl "comiclaser/model"
	"comiclaser/monitor"
	"comiclaser/setting"
	"path/filepath"

	"net/http"
	"net/url"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

func init() {
	setting.Initialize()
	ll, err := log.ParseLevel(setting.LogLevel)
	if err == nil {
		log.SetLevel(ll)
	}
}

func main() {
	mdl.SetupDb(setting.DatabaseName)

	monitor.Remove(setting.ComicFolder)
	monitor.Update(setting.ComicFolder)
	go monitor.Watch(setting.ComicFolder)

	e := echo.New()

	if setting.ServerLogging {
		e.Use(serverLogger())
	}

	e.Get("/", root)
	e.Get("/dbinfo", dbinfo)
	e.Get("/comic/:id", getComic)
	e.Get("/comic/:id/page/:page", getPage)
	e.Get("/comiclist", comicList)
	e.Get("/folders", folders)
	e.Get("/folders/*", subFolders)
	log.Info("Server (version " + setting.APP_VER + ") started on port " + setting.HttpPort)
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
		len(*comics),
		len(*comics),
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

func comicList(c *echo.Context) error {
	f := c.Query("folder")
	var comics *[]mdl.Comic
	var err error
	if f == "" {
		f = setting.ComicFolder
		comics, err = mdl.GetAllComicsInFolder(f)
	} else {
		comics, err = mdl.GetAllComics()
	}

	if err != nil {
		log.WithError(err).Error("Problem getting all comics")
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to get all comics")
	}
	data := struct {
		Comics *[]mdl.Comic `json:"comics"`
		Total  int          `json:"total_count"`
		Pages  int          `json:"page_count"`
	}{
		comics,
		len(*comics),
		len(*comics),
	}
	return c.JSON(http.StatusOK, data)

}

func folders(c *echo.Context) error {
	folders, comicCount := getStructure(setting.ComicFolder)

	comicListEndpoint := ""
	if comicCount > 0 {
		v := url.Values{}
		v.Add("folder", setting.ComicFolder)

		comicListEndpoint = "/comiclist?" + v.Encode()
	}

	comics := struct {
		Count int    `json:"count"`
		URL   string `json:"url_path"`
	}{
		comicCount,
		comicListEndpoint,
	}

	data := struct {
		Current string      `json:"current"`
		Folders []folder    `json:"folders"`
		Comics  interface{} `json:"comics"`
	}{
		"",
		folders,
		comics,
	}
	return c.JSON(http.StatusOK, data)
}

func subFolders(c *echo.Context) error {
	n := c.Param("_*")

	subFolders, comicCount := getSubStructure(setting.ComicFolder, n)

	comicListEndpoint := ""
	if comicCount > 0 {
		v := url.Values{}
		v.Add("folder", filepath.Join(setting.ComicFolder, n))

		comicListEndpoint = "/comiclist?" + v.Encode()
	}

	comics := struct {
		Count int    `json:"count"`
		URL   string `json:"url_path"`
	}{
		comicCount,
		comicListEndpoint,
	}

	data := struct {
		Current string      `json:"current"`
		Folders []folder    `json:"folders"`
		Comics  interface{} `json:"comics"`
	}{
		filepath.Join(setting.ComicFolder, n),
		subFolders,
		comics,
	}

	return c.JSON(http.StatusOK, data)
}
