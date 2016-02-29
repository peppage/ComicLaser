package model

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // required for sqlx

	log "github.com/Sirupsen/logrus"
)

var db *sqlx.DB

func init() {
	var err error
	db, err = sqlx.Connect("sqlite3", "./comicdb2.sqlite")
	if err != nil {
		log.WithError(err).Panic("Cannot connect to database")
	}
}

// SetupDb first time setup, creates tables if required.
func SetupDb() {
	tx := db.MustBegin()

	tx.Exec(`CREATE TABLE IF NOT EXISTS comics (path TEXT PRIMARY KEY, filename TEXT, series TEXT,
            size INTEGER, pages INTEGER, issue INTEGER)`)

	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		log.WithError(err).Error("DB setup failed.")
	}
}
