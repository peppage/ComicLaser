package model

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // required for sqlx

	log "github.com/Sirupsen/logrus"
)

var db *sqlx.DB

// SetupDb first time setup, creates tables if required.
func SetupDb(dbName string) {

	var err error
	db, err = sqlx.Connect("sqlite3", "./"+dbName)
	if err != nil {
		log.WithError(err).Panic("Cannot connect to database")
	}

	tx := db.MustBegin()

	tx.Exec(`CREATE TABLE IF NOT EXISTS comics (path TEXT PRIMARY KEY, filename TEXT, series TEXT,
            size INTEGER, pages INTEGER, issue INTEGER)`)
	tx.Exec(`CREATE TABLE IF NOT EXISTS dbinfo (name text PRIMARY KEY, value text)`)

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.WithError(err).Error("DB setup failed.")
	}

	var created int
	db.Get(&created, `SELECT value FROM dbinfo WHERE name=$1`, "created")
	if created == 0 {
		db.Exec(`INSERT INTO dbinfo (name, value) VALUES ($1, $2)`, "created", time.Now().Unix())
		db.Exec(`INSERT INTO dbinfo (name, value) VALUES ($1, $2)`, "last_update", time.Now().Unix())
	}
}
