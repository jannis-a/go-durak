package auth_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/raja/argon2pw"
	"github.com/stretchr/testify/assert"

	"github.com/jannis-a/go-durak/app"
	"github.com/jannis-a/go-durak/auth"
	"github.com/jannis-a/go-durak/utils"
)

var (
	a        *app.App
	id       uint
	username string
	password string
)

func setUp() {
	username = randomdata.SillyName()
	password = randomdata.RandStringRunes(randomdata.Number(8, 32))

	hashed, err := argon2pw.GenerateSaltedHash(password)
	if err != nil {
		log.Panic(err)
	}

	qry := `insert into users (username, email, password) values ($1, $2, $3) returning id`
	res := a.DB.QueryRow(qry, username, randomdata.Email(), hashed)
	if err := res.Scan(&id); err != nil {
		log.Panic(err)
	}
}

func tearDown() {
	a.DB.Exec(`truncate table users, tokens cascade`)
}

func TestMain(m *testing.M) {
	a = app.NewApp()
	a.RegisterApi("auth", auth.Routes)

	setUp()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

func doLogin(t *testing.T, data map[string]string) *httptest.ResponseRecorder {
	payload, err := json.Marshal(data)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(payload))
	assert.Nil(t, err)

	return utils.DispatchRequest(a.Router, req)
}

func TestLoginInvalidPayload(t *testing.T) {
	data := []map[string]string{
		{},
		{"username": "user"},
		{"password": "secret"},
	}

	for _, d := range data {
		res := doLogin(t, d)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	res := doLogin(t, map[string]string{
		"username": username,
		"password": "INVALID",
	})

	assert.Equal(t, http.StatusUnauthorized, res.Code)
}

func TestLoginValidCredentials(t *testing.T) {
	res := doLogin(t, map[string]string{
		"username": username,
		"password": password,
	})

	assert.Equal(t, http.StatusOK, res.Code)
	assert.NotEmpty(t, res.Body.String())

	var cookie http.Cookie
	for _, c := range res.Result().Cookies() {
		if c.Name == auth.RefreshCookieName {
			cookie = *c
			break
		}
	}

	assert.NotNil(t, cookie)
	assert.True(t, cookie.HttpOnly)
	assert.NotEmpty(t, cookie.Value)

	var count int
	qry := `select count(*) from tokens where token = $1 and user_id = $2`
	row := a.DB.QueryRow(qry, cookie.Value, id)
	assert.Nil(t, row.Scan(&count))
	assert.Equal(t, 1, count)
}

func TestRefreshMissingCookie(t *testing.T) {
	t.Fatal("Not implemented")
}

func TestRefreshInvalidToken(t *testing.T) {
	t.Fatal("Not implemented")
}

func TestRefreshValidToken(t *testing.T) {
	t.Fatal("Not implemented")
}

func TestLogoutMissingCookie(t *testing.T) {
	t.Fatal("Not implemented")
}

func TestLogoutInvalidToken(t *testing.T) {
	t.Fatal("Not implemented")
}

func TestLogoutValidToken(t *testing.T) {
	t.Fatal("Not implemented")
}
