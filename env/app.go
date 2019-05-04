package env

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Router *mux.Router
	DB     *sqlx.DB
	Config *Config
}

func NewApp(config *Config) *App {
	if config == nil {
		config = NewConfig()
	}

	return &App{
		Config: config,
		Router: mux.NewRouter().StrictSlash(true),
		DB:     NewDatabase(config),
	}
}
