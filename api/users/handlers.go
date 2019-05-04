package users

import (
	"net/http"
	"net/url"

	"github.com/go-chi/render"
	"github.com/thedevsaddam/govalidator"

	"github.com/jannis-a/go-durak/env"
	"github.com/jannis-a/go-durak/handler"
	"github.com/jannis-a/go-durak/utils"
)

func ListHandler(a *env.App, w http.ResponseWriter, r *http.Request) error {
	users := make([]User, 0)

	rows, _ := a.DB.Queryx(`SELECT * FROM users`)
	for rows.Next() {
		var u User
		err := rows.StructScan(&u)
		if err == nil {
			users = append(users, u)
		}
	}

	render.JSON(w, r, users)
	return nil
}

func CreateHandler(a *env.App, w http.ResponseWriter, r *http.Request) error {
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

	row := a.DB.QueryRow(`SELECT * FROM users WHERE username = $1`, d["username"])
	if _ = row.Scan(count); 0 < count {
		e.Add("username", "USERNAME ALREADY TAKEN")
	}

	row = a.DB.QueryRow(`SELECT * FROM users WHERE email = $1`, d["email"])
	if _ = row.Scan(count); 0 < count {
		e.Add("email", "EMAIL ALREADY TAKEN")
	}

	if 0 < len(e) {
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, map[string]url.Values{"errors": e})
		return nil
	}

	user := NewUser(a.DB, d["username"], d["email"], d["password"])
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)

	return nil
}

func DetailHandler(a *env.App, w http.ResponseWriter, r *http.Request) error {
	username := utils.GetRouteParam(r, "username")

	// language=SQL
	row := a.DB.QueryRowx(`SELECT * FROM users WHERE username = $1`, username)

	var user User
	err := row.StructScan(&user)
	if err != nil {
		return handler.NewStatusError(http.StatusNotFound, "")
	}

	render.JSON(w, r, user)
	return nil
}

func UpdateHandler(a *env.App, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func DeleteHandler(a *env.App, w http.ResponseWriter, r *http.Request) error {
	username := utils.GetRouteParam(r, "username")

	result, _ := a.DB.Exec(`DELETE FROM users WHERE username = $1`, username)
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return handler.NewStatusError(http.StatusNotFound, "")
	}

	render.Status(r, http.StatusAccepted)
	render.PlainText(w, r, "")
	return nil
}
