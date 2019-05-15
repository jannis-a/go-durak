package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"

	"github.com/jannis-a/go-durak/app"
	"github.com/jannis-a/go-durak/users"
	"github.com/jannis-a/go-durak/utils"
)

var a *app.App

func truncateTables() {
	a.DB.Exec(`truncate table users cascade;
                   alter sequence users_id_seq restart;`)
}

func setUp(t *testing.T) func(*testing.T) {
	t.Log("Setup tables")

	return func(t *testing.T) {
		t.Log("Teardown tables")

		truncateTables()
	}
}

func createUser() users.User {
	return users.New(a, randomdata.SillyName(), randomdata.Email(), "secret")
}

func createUserPub() users.UserPub {
	user := createUser()
	return users.UserPub{
		Id:       user.Id,
		Username: user.Username,
		JoinedAt: user.JoinedAt,
	}
}

func TestMain(m *testing.M) {
	a = app.NewApp()
	a.RegisterApi("users", users.Routes)

	truncateTables()
	os.Exit(m.Run())
}

func TestList(t *testing.T) {
	testCases := []struct {
		want int
		page int
		size int
	}{
		{10, 0, 0},
		{10, 1, 0},
		{5, 3, 0},
		{8, 0, 8},
		{1, 4, 8},
	}

	tearDown := setUp(t)
	defer tearDown(t)

	expected := make([]users.UserPub, 0)
	for i := 0; i < 25; i++ {
		expected = append(expected, createUserPub())
	}

	for _, tc := range testCases {
		tcName := fmt.Sprintf("want=%d_page=%d_size=%d", tc.want, tc.page, tc.page)
		t.Run(tcName, func(t *testing.T) {

			params := make(url.Values)
			if 0 < tc.page {
				params.Add("page", strconv.Itoa(tc.page))
			}

			if 0 < tc.size {
				params.Add("size", strconv.Itoa(tc.size))
			}

			req, err := http.NewRequest("GET", "/users?"+params.Encode(), nil)
			assert.Nil(t, err)

			res := utils.DispatchRequest(a.Router, req)
			assert.Equal(t, http.StatusOK, res.Code)

			var result []users.UserPub
			err = json.Unmarshal(res.Body.Bytes(), &result)

			assert.Nil(t, err)
			assert.Len(t, result, tc.want)
		})
	}
}

func TestDetail(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown(t)

	expected := createUserPub()

	req, err := http.NewRequest("GET", "/users/"+expected.Username, nil)
	assert.Nil(t, err)

	res := utils.DispatchRequest(a.Router, req)
	assert.Equal(t, http.StatusOK, res.Code)

	var result users.UserPub
	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.Nil(t, err)

	assert.Equal(t, expected.Id, result.Id)
	assert.Equal(t, expected.Username, result.Username)
	assert.True(t, expected.JoinedAt.Equal(result.JoinedAt))
}

func TestCreate(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown(t)

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

	var result users.User
	row := a.DB.QueryRow(`select * from users where username = $1`, data["username"])
	err = row.Scan(&result.Id, &result.Username, &result.Email, &result.Password, &result.JoinedAt)
	assert.Nil(t, err)

	assert.Equal(t, data["username"], result.Username)
	assert.Equal(t, data["email"], result.Email)

	verified, err := utils.Argon2Verify(data["password"], result.Password)
	assert.Nil(t, err)
	assert.True(t, verified)
}
