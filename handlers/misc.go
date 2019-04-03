package handlers

import "github.com/jinzhu/gorm"

type App struct {
	db *gorm.DB
}

func NewApp(db *gorm.DB) App {
	return App{
		db: db,
	}
}
