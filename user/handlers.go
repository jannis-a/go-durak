package user

import (
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/raja/argon2pw"
	"github.com/thedevsaddam/govalidator"

	"github.com/jannis-a/go-durak/config"
)

func ListHandler(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []User
		c.Db.Find(&users)

		render.JSON(w, r, users)
	}
}

func CreateHandler(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := make(map[string]string)
		v := govalidator.New(govalidator.Options{
			Request:         r,
			Data:            &d,
			RequiredDefault: true,
			Rules: govalidator.MapData{
				"username":         []string{"required", "min:3", "max:50"},
				"email":            []string{"required", "email"},
				"password":         []string{"required"},
				"password_confirm": []string{"required"},
			},
		})
		e := v.ValidateJSON()

		if d["password"] != d["password_confirm"] {
			e.Add("password", "PW MISSMATCH")
		}

		var count int
		var user User

		c.Db.Model(&user).Where("username = ?", d["username"]).Count(&count)
		if 0 < count {
			e.Add("username", "USERNAME ALREADY TAKEN")
		}

		c.Db.Model(&user).Where("email = ?", d["email"]).Count(&count)
		if 0 < count {
			e.Add("email", "EMAIL ALREADY TAKEN")
		}

		if 0 < len(e) {
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, map[string]url.Values{"errors": e})
			return
		}

		hashedPassword, err := argon2pw.GenerateSaltedHash(d["password"])
		if err != nil {
			log.Panicf("Hash generated returned error: %v", err)
		}

		user = User{
			Username: d["username"],
			Email:    d["email"],
			Password: hashedPassword,
			Token:    generateToken(c.Db),
		}
		c.Db.Create(&user)

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, user)
	}
}

func DetailHandler(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		var user User
		result := c.Db.Where("username = ?", username).First(&user)

		if 0 == result.RowsAffected {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, "")
		} else {
			render.JSON(w, r, user)
		}
	}
}

func UpdateHandler(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func DeleteHandler(c *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		var user User
		result := c.Db.Where("username = ?", username).Delete(&user)

		if 0 == result.RowsAffected {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, "")
		} else {
			render.Status(r, http.StatusAccepted)
			render.PlainText(w, r, "")
		}
	}
}
