package users_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"

	"github.com/jannis-a/go-durak/api/users"
	"github.com/jannis-a/go-durak/env"
	"github.com/jannis-a/go-durak/utils"
)

var app *env.App

func truncateTable() {
	_, _ = app.DB.Query(`TRUNCATE TABLE users`)
}

func TestMain(m *testing.M) {
	app = env.NewApp(nil)
	users.Initialize(app)

	code := m.Run()
	truncateTable()
	os.Exit(code)
}

func TestList(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	assert.Nil(t, err)

	res := utils.DispatchRequest(app.Router, req)
	assert.Equal(t, res.Code, http.StatusOK)
	assert.Equal(t, res.Body.String(), "[]\n")
}

func TestDetail(t *testing.T) {
	user := users.NewUser(app.DB, randomdata.SillyName(), randomdata.Email(), "secret")
	expected := fmt.Sprintf(`{"username":"%s","email":"%s","joined_at":"%s"}`,
		user.Username,
		user.Email,
		user.JoinedAt.Format(time.RFC3339Nano),
	)

	req, err := http.NewRequest("GET", "/users/"+user.Username, nil)
	assert.Nil(t, err)

	res := utils.DispatchRequest(app.Router, req)
	assert.Equal(t, res.Code, http.StatusOK)
	assert.Equal(t, res.Body.String(), expected+"\n")
}
