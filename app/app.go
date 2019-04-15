package app

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	HOST string
	PORT string

	Postgres struct {
		HOST     string
		PORT     int
		SSL      string
		NAME     string
		USER     string
		PASSWORD string
	}
}

type App struct {
	Config
	Db *sqlx.DB
}

func New() (*App, error) {
	app := App{}
	cfg, err := initViper()
	app.Config = cfg

	if err != nil {
		return &app, err
	}

	database, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Postgres.HOST,
			cfg.Postgres.PORT,
			cfg.Postgres.USER,
			cfg.Postgres.PASSWORD,
			cfg.Postgres.NAME,
			cfg.Postgres.SSL,
		),
	)

	if err != nil {
		return &app, err
	}

	app.Db = database
	return &app, err
}

func NewTesting() (*App, error) {
	viper.AddConfigPath("..")
	return New()
}

func initViper() (Config, error) {
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		return Config{}, err
	}

	var constants Config
	err = viper.Unmarshal(&constants)
	return constants, nil
}
