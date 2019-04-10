package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

type Constants struct {
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

type Config struct {
	Constants
	Db *gorm.DB
}

func New() (*Config, error) {
	config := Config{}
	constants, err := initViper()
	config.Constants = constants

	if err != nil {
		return &config, err
	}

	database, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			constants.Postgres.HOST,
			constants.Postgres.PORT,
			constants.Postgres.USER,
			constants.Postgres.PASSWORD,
			constants.Postgres.NAME,
			constants.Postgres.SSL,
		),
	)

	if err != nil {
		return &config, err
	}

	config.Db = database
	return &config, err
}

func initViper() (Constants, error) {
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		return Constants{}, err
	}

	var constants Constants
	err = viper.Unmarshal(&constants)
	return constants, nil
}
