package user

import (
	"time"
)

type User struct {
	Id        uint      `gorm:"primary_key" json:"-"`
	Username  string    `gorm:"not_null;unique" json:"username"`
	Email     string    `gorm:"not_null;unique" json:"email"`
	Password  string    `gorm:"not_null" json:"-"`
	CreatedAt time.Time `gorm:"column:joined_at" json:"joined_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
