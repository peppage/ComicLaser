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
		log.WithField("error", err).Panic("Cannot connect to database")
	}
}

func SetupDb() {

}
