package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var db *gorm.DB

func Open() *gorm.DB {
	var err error

	db, err = gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			viper.GetString("DB_HOST"),
			viper.GetInt("DB_PORT"),
			viper.GetString("DB_USER"),
			viper.GetString("DB_PASSWORD"),
			viper.GetString("DB_NAME"),
			viper.GetString("DB_SSL"),
		),
	)

	if err != nil {
		panic(fmt.Errorf("Error connecting to database: %s \n", err))
	}

	return db
}
