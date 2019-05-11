package app

import (
	"database/sql"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func NewDatabase(config *Config) *sql.DB {
	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		log.Panic(err)
	}

	return db
}
