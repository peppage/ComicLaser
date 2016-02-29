package model

type Comic struct {
	Path     string
	FileName string
	Size     int64
	Series   string
	Pages    int
	Issue    int
}

}

func CreateComic(c Comic) error {

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
