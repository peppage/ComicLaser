package model

import (
	"bytes"
	"io"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"comiclaser/lzmadec"
)

type Comic struct {
	ID       int64  `json:"id"`
	Path     string `json:"path"`
	FileName string `json:"filaname"`
	Folder   string `json:"folder"`
	Size     int64  `json:"filesize"`
	Series   string `json:"series"`
	Pages    int    `json:"page_count"`
	Issue    int    `json:"issue"`
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
	tx.Exec(`INSERT INTO comics (path, filename, folder, series, size, pages, issue) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		c.Path, c.FileName, c.Folder, c.Series, c.Size, c.Pages, c.Issue)
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

func GetComic(id int) (*[]Comic, error) {
	c := []Comic{}

	err := db.Select(&c, `SELECT * FROM comics WHERE id=$1`, id)

	return &c, err
}

func GetComicPage(id int, page int) ([]byte, error) {
	c := Comic{}
	err := db.Get(&c, `SELECT * FROM comics WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}

	a, err := lzmadec.NewArchive(c.Path)
	if err != nil {
		return nil, err
	}

	sort.Sort(lzmadec.ByPath(a.Entries))

	r, err := a.GetFileReader(a.Entries[page].Path)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	r.Close()
	return buf.Bytes(), nil

}

func GetAllComics(folder string) (*[]Comic, error) {
	c := []Comic{}
	err := db.Select(&c, `SELECT * FROM comics WHERE folder=$1`, folder)
	return &c, err
}

func RemoveComic(id int64) error {
	_, err := db.Exec(`DELETE FROM comics WHERE id=$1`, id)
	return err
}
