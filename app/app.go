package app

import (
	"database/sql"

	"github.com/gorilla/mux"

	"github.com/jannis-a/go-durak/utils"
)

type App struct {
	Router       *mux.Router
	DB           *sql.DB
	Config       *Config
	Argon2Params *utils.Argon2Params
}

func NewApp() *App {
	config := NewConfig()

	return &App{
		Config: config,
		Router: mux.NewRouter().StrictSlash(true),
		DB:     NewDatabase(config),
		Argon2Params: &utils.Argon2Params{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
}
