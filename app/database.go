package app

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewDatabase(config *Config) *sql.DB {
	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		log.Panic(err)
	}

	return db
}
