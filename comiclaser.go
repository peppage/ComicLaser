package main

import (
	"bufio"
	"bytes"
	"errors"
	"strings"

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

	/*cmd := exec.Command("7z", "l", "-slt", "comics\\A-Force 003 (2015) (Digital) (Zone-Empire).cbr")
	out, err := cmd.CombinedOutput()
	log.WithFields(log.Fields{
		"cmd": cmd,
		//	"out": out,
		"err": err,
	}).Debug("opening file")
	parse7zListOutput(out)*/

	go monitor.Watch(setting.ComicFolder)

	e := echo.New()
	e.Get("/", root)
	e.Run(":8000")

}

func root(c *echo.Context) error {
	return c.String(http.StatusOK, "st")
}

func parse7zListOutput(d []byte) {
	//var res []Entry
	r := bytes.NewBuffer(d)
	scanner := bufio.NewScanner(r)
	err := advanceToFirstEntry(scanner)
	if err != nil {
		log.Error(err)
	}
	for {
		lines, err := getEntryLines(scanner)
		if err != nil {
			log.Error(err)
		}
		if len(lines) == 0 {
			// last entry
			break
		}
		/*e, err := parseEntryLines(lines)
		if err != nil {
			return nil, err
		}
		res = append(res, e)*/
	}
	//return res, nil
}

func getEntryLines(scanner *bufio.Scanner) ([]string, error) {
	var res []string
	for scanner.Scan() {
		s := scanner.Text()
		s = strings.TrimSpace(s)
		if s == "" {
			break
		}
		log.Debug(s)
		res = append(res, s)
	}
	err := scanner.Err()
	if err != nil {
		return nil, err
	}
	if len(res) == 9 || len(res) == 0 {
		return res, nil
	}
	return nil, errors.New("too many lines")
}

func advanceToFirstEntry(scanner *bufio.Scanner) error {
	for scanner.Scan() {
		s := scanner.Text()
		if s == "----------" {
			return nil
		}
	}
	err := scanner.Err()
	if err == nil {
		err = errors.New("no entries")
	}
	return err
}
