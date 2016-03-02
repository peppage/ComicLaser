package model

import "time"

type DbInfo struct {
	Created     string `json:"created"`
	Updated     string `json:"last_updated"`
	TotalComics int    `json:"comic_count"`
	ID          string `json:"id"`
}

// DbUpdated sets that this is the time the database was last updated
func DbUpdated() error {
	_, err := db.Exec(`UPDATE dbinfo SET value=$1 WHERE name = $2`, time.Now().Format("2006-01-02T15:04:05.999999"), "last_update")
	return err
}

func GetDbInfo() DbInfo {
	var i DbInfo
	db.Get(&i.Created, `SELECT value FROM dbinfo WHERE name=$1`, "created")
	db.Get(&i.Updated, `SELECT value FROM dbinfo WHERE name=$1`, "last_update")
	db.Get(&i.TotalComics, `SELECT COUNT(*) FROM comics as totalcomics`)
	db.Get(&i.ID, `SELECT value FROM dbinfo WHERE name=$1`, "id")
	return i
}
