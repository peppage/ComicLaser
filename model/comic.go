package model

import (
	"path/filepath"
	"strings"

	"ComicLaser/lzmadec"

	log "github.com/Sirupsen/logrus"
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
	log.Debug(c.Pages)
	if c.Pages > 0 {
		pn := parseFileName(a.Entries[0].Path)
		c.Series = pn.Series
	}

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
	s := strings.Split(file, "(")
	if len(s) > 0 {
		pn.Series = strings.TrimSpace(s[0])
	}

	return pn
}
