package users

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/thedevsaddam/govalidator"

	"github.com/jannis-a/go-durak/app"
)

func ListHandler(c *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []User

		rows, _ := c.Db.Queryx(`SELECT * FROM users`)
		for rows.Next() {
			var u User
			err := rows.StructScan(&u)
			if err == nil {
				users = append(users, u)
			}
		}

		render.JSON(w, r, users)
	}
}

func CreateHandler(c *app.App) http.HandlerFunc {
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
		// var user User

		row := c.Db.QueryRow(`SELECT * FROM users WHERE username = ?`, d["username"])
		if _ = row.Scan(count); 0 < count {
			e.Add("username", "USERNAME ALREADY TAKEN")
		}

		row = c.Db.QueryRow(`SELECT * FROM users WHERE email = ?`, d["email"])
		if _ = row.Scan(count); 0 < count {
			e.Add("email", "EMAIL ALREADY TAKEN")
		}

		if 0 < len(e) {
			render.Status(r, http.StatusUnprocessableEntity)
			render.JSON(w, r, map[string]url.Values{"errors": e})
			return
		}

		user := NewUser(c.Db, d["username"], d["email"], d["password"])

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, user)
	}
}

func DetailHandler(c *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		var user User
		// language=SQL
		row := c.Db.QueryRowx(`SELECT * FROM users WHERE username = $1`, username)
		_ = row.StructScan(&user)

		if &user == nil {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, "")
		} else {
			render.JSON(w, r, user)
		}
	}
}

func UpdateHandler(c *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func DeleteHandler(c *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		result, _ := c.Db.Exec(`DELETE FROM users WHERE username = $1`, username)
		if rows, _ := result.RowsAffected(); rows == 0 {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, "")
		} else {
			render.Status(r, http.StatusAccepted)
			render.PlainText(w, r, "")
		}
	}
}
