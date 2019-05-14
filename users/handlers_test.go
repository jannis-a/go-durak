package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"

	"github.com/jannis-a/go-durak/app"
	"github.com/jannis-a/go-durak/users"
	"github.com/jannis-a/go-durak/utils"
)

var a *app.App

func truncateTable() {
	_, _ = a.DB.Exec(`truncate table users cascade;
                          alter sequence users_id_seq restart;`)
}

func createUser() users.User {
	return users.New(a, randomdata.SillyName(), randomdata.Email(), "secret")
}

func createUserPub() users.UserPub {
	user := createUser()
	return users.UserPub{user.Id, user.Username, user.JoinedAt}
}

func TestMain(m *testing.M) {
	a = app.NewApp()
	a.RegisterApi("users", users.Routes)

	truncateTable()
	code := m.Run()
	truncateTable()
	os.Exit(code)
}

func TestList(t *testing.T) {
	expected := []users.UserPub{createUserPub()}

	req, err := http.NewRequest("GET", "/users", nil)
	assert.Nil(t, err)

	res := utils.DispatchRequest(a.Router, req)
	assert.Equal(t, http.StatusOK, res.Code)

	var result []users.UserPub
	err = json.Unmarshal(res.Body.Bytes(), &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestDetail(t *testing.T) {
	expected := createUserPub()

	req, err := http.NewRequest("GET", "/users/"+expected.Username, nil)
	assert.Nil(t, err)

	res := utils.DispatchRequest(a.Router, req)
	assert.Equal(t, http.StatusOK, res.Code)

	var result users.UserPub
	err = json.Unmarshal(res.Body.Bytes(), &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestCreate(t *testing.T) {
	data := map[string]string{
		"username":         randomdata.SillyName(),
		"email":            randomdata.Email(),
		"password":         "secret",
		"password_confirm": "secret",
	}
	payload, err := json.Marshal(data)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
	assert.Nil(t, err)

	res := utils.DispatchRequest(a.Router, req)
	assert.Equal(t, http.StatusCreated, res.Code)

	row := a.DB.QueryRow(`select * from users where username = $1`, data["username"])
	var user users.User
	err = row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.JoinedAt)
	assert.Nil(t, err)

	expected := fmt.Sprintf(`{"id":%d,"username":"%s","joined_at":"%s","email":"%s"}`,
		user.Id,
		user.Username,
		user.JoinedAt.Format(time.RFC3339Nano),
		user.Email,
	)
	assert.Equal(t, data["username"], user.Username)
	assert.Equal(t, data["email"], user.Email)
	result, err := utils.Argon2Verify(data["password"], user.Password)
	assert.Nil(t, err)
	assert.True(t, result)
	assert.Equal(t, expected, res.Body.String())
}
