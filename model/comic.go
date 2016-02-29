package model

type Comic struct {
	ID     int
	Path   string
	Series string
}

func CreateComic(c Comic) error {

	tx := db.MustBegin()
	tx.Exec(`INSERT INTO comics (path, series) VALUES ($1, $2)`, c.Path, c.Series)
	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
