package users

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/jannis-a/go-durak/app"
	"github.com/jannis-a/go-durak/utils"
)

func ListHandler(a *app.App, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	size, err := strconv.Atoi(query.Get("size"))
	if err != nil {
		size = 10
	}

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 1
	}

	qry := `select id, username, joined_at from users limit $1 offset $2`
	rows, err := a.DB.Query(qry, size, size*(page-1))
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
	var (
		data           UserCreate
		usernameExists bool
		emailExists    bool
	)

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		utils.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	e := make(url.Values)

	qry := `select
		max(case when username = $1 then 1 else 0 end),
	  max(case when email = $2 then 1 else 0 end)
	from users`
	row := a.DB.QueryRow(qry, data.Username, data.Email)
	_ = row.Scan(&usernameExists, &emailExists)

	if 3 > len(data.Username) || 50 < len(data.Username) {
		e.Add("username", "Username length must be between 3 and 50")
	}

	// TODO: assert email valid

	if data.Password != data.PasswordConfirm {
		e.Add("password", "Passwords don't match")
	}

	// TODO: password rules

	if usernameExists {
		e.Add("username", "Username already taken")
	}

	if emailExists {
		e.Add("email", "Email already taken")
	}

	if 0 < len(e) {
		utils.RenderErrors(w, e)
		return
	}

	user := New(a, data.Username, data.Email, data.Password)

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
