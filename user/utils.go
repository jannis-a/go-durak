package user

import (
	"github.com/jinzhu/gorm"

	"github.com/jannis-a/go-durak/utils"
)

func generateToken(db *gorm.DB) string {
	var token string
	var result int

	for {
		token = utils.RandString(32)
		db.Table("users").Where("token = ?", token).Count(&result)

		if 0 == result {
			return token
		}
	}
}
