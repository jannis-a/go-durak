package env

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDatabase(config *Config) *sqlx.DB {
	db, err := sqlx.Open("postgres", config.DB)
	if err != nil {
		log.Panic(err)
	}

	return db
}
