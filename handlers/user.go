package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jannis-a/go-durak/models"
	"github.com/jannis-a/go-durak/utils"
	"github.com/raja/argon2pw"
	"github.com/thedevsaddam/govalidator"
	"log"
	"net/http"
)

func (a *App) generateToken() string {
	var token string
	var result int

	for {
		token = utils.RandString(32)
		a.db.Table("users").Where("token = ?", token).Count(&result)

		if 0 == result {
			return token
		}
	}
}

func (a *App) UserList(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	a.db.Find(&users)

	render.JSON(w, r, users)
}

func (a *App) UserCreate(w http.ResponseWriter, r *http.Request) {
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

	exists := func(key string, value string) bool {
		var count int
		query := fmt.Sprintf("%s = ?", key)
		a.db.Model(&models.User{}).Where(query, value).Count(&count)

		return 0 < count
	}

	if exists("username", d["username"]) {
		print("NAME EXISTS")
		e.Add("username", "USERNAME ALREADY TAKEN")
	}

	if exists("email", d["email"]) {
		print("EMAIL EXISTS")
		e.Add("email", "EMAIL ALREADY TAKEN")
	}

	if 0 < len(e) {
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, map[string]interface{}{"errors": e})
		return
	}

	hashedPassword, err := argon2pw.GenerateSaltedHash(d["password"])
	if err != nil {
		log.Panicf("Hash generated returned error: %v", err)
	}

	user := models.User{
		Username: d["username"],
		Email:    d["email"],
		Password: hashedPassword,
		Token:    a.generateToken(),
	}
	a.db.Create(&user)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}

func (a *App) findOr404(
	w http.ResponseWriter,
	r *http.Request,
	handler func(w http.ResponseWriter, r *http.Request, user models.User),
) {
	var user models.User
	query := a.db.Where("username = ?", chi.URLParam(r, "username")).First(&user)

	if 0 == query.RowsAffected {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, "")
		return
	}

	handler(w, r, user)
}

func (a *App) UserDetail(w http.ResponseWriter, r *http.Request) {
	a.findOr404(w, r, func(w http.ResponseWriter, r *http.Request, user models.User) {
		render.JSON(w, r, user)
	})
}

func (a *App) UserUpdate(w http.ResponseWriter, r *http.Request) {
	a.findOr404(w, r, func(w http.ResponseWriter, r *http.Request, user models.User) {
		// TODO: make update

		render.JSON(w, r, user)
	})
}

func (a *App) UserDelete(w http.ResponseWriter, r *http.Request) {
	a.findOr404(w, r, func(w http.ResponseWriter, r *http.Request, user models.User) {
		a.db.Delete(user)

		render.Status(r, http.StatusAccepted)
		render.PlainText(w, r, "")
	})
}
