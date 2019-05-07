package auth_test

import (
	"os"
	"testing"

	"github.com/jannis-a/go-durak/app"
	"github.com/jannis-a/go-durak/auth"
)

var a *app.App

func tearDown() {
	_, _ = a.DB.Exec(`truncate table users, tokens cascade`)
}

func TestMain(m *testing.M) {
	a = app.NewApp(nil)
	a.Register("auth", auth.Routes)
	// env.RegisterOld(app, "auth", auth.Routes)

	retCode := m.Run()
	tearDown()

	os.Exit(retCode)
}

func TestLoginInvalidPayload(t *testing.T) {
	t.Fatal("Not implemented")
}

func TestLoginInvalidCredentials(t *testing.T) {
	t.Fatal("Not implemented")
}

func TestLoginValidCredentials(t *testing.T) {
	t.Fatal("Not implemented")
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
