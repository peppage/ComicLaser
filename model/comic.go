package model

import (
	"bytes"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"ComicLaser/lzmadec"
)

type Comic struct {
	Path     string
	FileName string
	Size     int64
	Series   string
	Pages    int
	Issue    int
}

// CreateComic given a path creates a comic struct.
func CreateComic(path string) (*Comic, error) {
	c := &Comic{Path: path}
	a, err := lzmadec.NewArchive(path)
	if err != nil {
		return nil, err
	}

	c.Pages = len(a.Entries)
	pn := parseFileName(path)
	c.Series = pn.Series
	c.Issue = pn.Issue

	return c, nil
}

// SaveComic saves a comic to the database.
func SaveComic(c *Comic) error {

	tx := db.MustBegin()
	tx.Exec(`INSERT INTO comics (path, filename, series, size, pages, issue) VALUES ($1, $2, $3, $4, $5, $6)`,
		c.Path, c.FileName, c.Series, c.Size, c.Pages, c.Issue)
	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

type parsedName struct {
	Series string
	Issue  int
}

func parseFileName(path string) parsedName {
	pn := parsedName{}
	_, file := filepath.Split(path)

	parens := regexp.MustCompile("\\(.*?\\)")

	f := string(parens.ReplaceAll([]byte(file), []byte("")))

	splits := strings.Split(f, " ")

	issueLoc := 0
	for i, v := range splits {
		if j, err := strconv.ParseInt(v, 10, 32); err == nil {
			pn.Issue = int(j)
			issueLoc = i
		}
	}

	var title bytes.Buffer
	for i := 0; i < issueLoc; i++ {
		title.WriteString(splits[i] + " ")
	}
	pn.Series = strings.TrimSpace(title.String())

	return pn
}
