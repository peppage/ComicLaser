package model

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // required for sqlx
	"github.com/satori/go.uuid"

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

	tx.Exec(`CREATE TABLE IF NOT EXISTS comics (id INTEGER PRIMARY KEY, path TEXT, filename TEXT UNIQUE, series TEXT,
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
		db.Exec(`INSERT INTO dbinfo (name, value) VALUES ($1, $2)`, "created", time.Now().Format("2006-01-02T15:04:05.999999"))
		db.Exec(`INSERT INTO dbinfo (name, value) VALUES ($1, $2)`, "last_update", time.Now().Format("2006-01-02T15:04:05.999999"))
		db.Exec(`insert INTO dbinfo (name, value) VALUES ($1, $2)`, "id", uuid.NewV4())
	}
}
