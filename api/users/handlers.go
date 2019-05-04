package users

import (
	"log"
	"net/http"
	"net/url"

	"github.com/thedevsaddam/govalidator"

	"github.com/jannis-a/go-durak/env"
	"github.com/jannis-a/go-durak/handler"
	"github.com/jannis-a/go-durak/utils"
)

func ListHandler(a *env.App, w http.ResponseWriter, r *http.Request) error {
	rows, err := a.DB.Query(`SELECT id, username, joined_at FROM users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var result []UserPub
	for rows.Next() {
		user := UserPub{}

		err := rows.Scan(&user.Id, &user.Username, &user.JoinedAt)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, user)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	utils.RenderJson(w, result)
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
		w.WriteHeader(http.StatusUnprocessableEntity)
		utils.RenderJson(w, map[string]url.Values{"errors": e})
		return nil
	}

	user := New(a.DB, d["username"], d["email"], d["password"])
	w.WriteHeader(http.StatusCreated)
	utils.RenderJson(w, user)
	return nil
}

func DetailHandler(a *env.App, w http.ResponseWriter, r *http.Request) error {
	username := utils.GetRouteParam(r, "username")

	row := a.DB.QueryRowx(`SELECT id, username, joined_at FROM users WHERE username = $1`, username)

	var user UserPub
	err := row.StructScan(&user)
	if nil != err {
		return handler.NewStatusError(http.StatusNotFound, "")
	}

	utils.RenderJson(w, user)
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

	w.WriteHeader(http.StatusAccepted)
	utils.RenderJson(w, nil)
	return nil
}
