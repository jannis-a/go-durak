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


func TestLoginHandler(t *testing.T) {
	testCases := []utils.ApiTestCase{{
		Name: "missing username and password",
		Data: map[string]string{},
		Code: http.StatusBadRequest,
	}, {
		Name: "missing password",
		Data: map[string]string{"username": username},
		Code: http.StatusBadRequest,
	}, {
		Name: "missing username",
		Data: map[string]string{"password": password},
		Code: http.StatusBadRequest,
	}, {
		Name: "invalid credentials",
		Data: map[string]string{"username": username, "password": "INVALID"},
		Code: http.StatusUnauthorized,
	}, {
		Name: "valid credentials",
		Data: map[string]string{"username": username, "password": password},
		Code: http.StatusOK,
		Func: func(t *testing.T, res *httptest.ResponseRecorder) {
			// Assert access token is not empty
			access := res.Body.String()
			assert.NotEmpty(t, access)

			// Check for refresh token cookie
			var refresh string
			for _, c := range res.Result().Cookies() {
				if c.Name == auth.RefreshCookieName {
					assert.NotEmpty(t, c.Value)
					assert.True(t, c.HttpOnly)
					refresh = c.Value
					break
				}
			}
			assert.NotEmpty(t, refresh)

			// Check for refresh token in database
			var count int
			qry := `select count(*) from tokens where user_id = $1 and token = $2`
			row := a.DB.QueryRow(qry, id, refresh)
			assert.Nil(t, row.Scan(&count))
			assert.Equal(t, 1, count)
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			data, err := json.Marshal(tc.Data)
			assert.Nil(t, err)

			req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(data))
			assert.Nil(t, err)

			res := utils.DispatchRequest(a.Router, req)
			assert.Equal(t, tc.Code, res.Code)

			if tc.Func != nil {
				tc.Func(t, res)
			}
		})
	}
}

func TestRefreshHandler(t *testing.T) {
	testCases := []utils.ApiTestCase{{
		Name: "missing cookie",
		Code: http.StatusBadRequest,
	}, {
		Name: "invalid token",
		Code: http.StatusUnauthorized,
	}, {
		Name: "valid token",
		Code: http.StatusOK,
	}}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Fatal("Not implemented")
		})
	}
}

func TestLogoutHandler(t *testing.T) {
	testCases := []utils.ApiTestCase{{
		Name: "missing cookie",
		Code: http.StatusBadRequest,
	}, {
		Name: "invalid token",
		Code: http.StatusUnauthorized,
	}, {
		Name: "valid token",
		Code: http.StatusOK,
	}}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Fatal("Not implemented")
		})
	}
}
