package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jannis-a/go-durak/app"
)

func main() {
	cfg := app.NewConfig()
	migrations := "./migrations"

	m, err := migrate.New("file://"+migrations, cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	_ = m.Up()
}
