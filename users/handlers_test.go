package users

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"

	"github.com/jannis-a/go-durak/app"
)

var c *app.App

func TestMain(m *testing.M) {
	var err error
	c, err = app.NewTesting()
	if err != nil {
		log.Panicln("Configuration error", err)
	}

	_, _ = c.Db.Query(`TRUNCATE TABLE users`)

	os.Exit(m.Run())
}

func TestListHandler(t *testing.T) {
	u := NewUser(c.Db, randomdata.SillyName(), randomdata.Email(), "secret")

	json := `[{"username":"%s","email":"%s","joined_at":"%s"}]` + "\n"
	expected := fmt.Sprintf(json, u.Username, u.Email, u.JoinedAt.Format(time.RFC3339Nano))

	req, err := http.NewRequest("GET", "/v1/users", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	ListHandler(c)(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, expected, w.Body.String())
}
