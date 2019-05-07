package app

import (
	"database/sql"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	Config *Config
}

func NewApp() *App {
	config := NewConfig()

	return &App{
		Config: config,
		Router: mux.NewRouter().StrictSlash(true),
		DB:     NewDatabase(config),
	}
}
