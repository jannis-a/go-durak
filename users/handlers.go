package users

import (
	"log"
	"net/http"
	"net/url"

	"github.com/thedevsaddam/govalidator"

	"github.com/jannis-a/go-durak/app"
	"github.com/jannis-a/go-durak/utils"
)

func ListHandler(a *app.App, w http.ResponseWriter, r *http.Request) {
	rows, err := a.DB.Query(`select id, username, joined_at from users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var result []UserPub
	for rows.Next() {
		var user UserPub
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
}

func CreateHandler(a *app.App, w http.ResponseWriter, r *http.Request) {
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

	// row := a.DB.QueryRow(`select
	// 	max(case when username = $1 then 1 else 0 end),
	//   max(case when email = $2 then 1 else 0 end)
	// from users`, d["username"], d["email"])

	row := a.DB.QueryRow(`select * from users where username = $1`, d["username"])
	if _ = row.Scan(count); 0 < count {
		e.Add("username", "USERNAME ALREADY TAKEN")
	}

	row = a.DB.QueryRow(`select * from users where email = $1`, d["email"])
	if _ = row.Scan(count); 0 < count {
		e.Add("email", "EMAIL ALREADY TAKEN")
	}

	if 0 < len(e) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		utils.RenderJson(w, map[string]url.Values{"errors": e})
		return
	}

	user := New(a.DB, d["username"], d["email"], d["password"])
	w.WriteHeader(http.StatusCreated)
	utils.RenderJson(w, user)
}

func DetailHandler(a *app.App, w http.ResponseWriter, r *http.Request) {
	var user UserPub
	qry := `select id, username, joined_at from users where username = $1`
	row := a.DB.QueryRow(qry, utils.GetRouteParam(r, "username"))
	err := row.Scan(&user.Id, &user.Username, &user.JoinedAt)
	if nil != err {
		utils.HttpError(w, http.StatusNotFound, "")
		return
	}

	utils.RenderJson(w, user)
}

func UpdateHandler(a *app.App, w http.ResponseWriter, r *http.Request) {
}

func DeleteHandler(a *app.App, w http.ResponseWriter, r *http.Request) {
	username := utils.GetRouteParam(r, "username")

	res, _ := a.DB.Exec(`delete from users where username = $1`, username)
	rows, _ := res.RowsAffected()
	if rows == 0 {
		utils.HttpError(w, http.StatusNotFound, "")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	utils.RenderJson(w, nil)
}
