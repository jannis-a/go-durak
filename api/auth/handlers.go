package auth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/raja/argon2pw"

	"github.com/jannis-a/go-durak/env"
	"github.com/jannis-a/go-durak/handler"
)

func LoginHandler(a *env.App, w http.ResponseWriter, r *http.Request) error {
	var (
		creds    Credentials
		subject  uint
		password string
	)

	// Decode credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		return handler.NewStatusError(http.StatusBadRequest, err.Error())
	}

	// Get user data and
	row := a.DB.QueryRow(`select id, password from users where username = $1`, creds.Username)
	if err := row.Scan(&subject, &password); err != nil {
		log.Fatal(err)
	}

	// Validate credentials
	valid, err := argon2pw.CompareHashWithPassword(password, creds.Password)
	if err != nil {
		return handler.NewStatusError(http.StatusInternalServerError, err.Error())
	} else if !valid {
		return handler.NewStatusError(http.StatusUnauthorized, "")
	}

	// Create refresh token
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	refresh := fmt.Sprintf("%x", bytes)

	// Insert refresh token into database
	_, err = a.DB.Exec(`insert into tokens (user_id, token) values ($1, $2)`, subject, refresh)
	if err != nil {
		log.Fatal(err)
	}

	// Create access token
	access := createAccessToken(a, subject, creds.Username)

	// Response
	http.SetCookie(w, &http.Cookie{
		Name:   RefreshCookieName,
		Value:  refresh,
		Secure: true,
	})
	w.Write([]byte(access))
	return nil
}

func RefreshHandler(a *env.App, w http.ResponseWriter, r *http.Request) error {
	var (
		subject  uint
		username string
	)

	// Get refresh token
	refresh, err := getRefreshToken(r)
	if err != nil {
		return err
	}

	// Validate token and fetch subject data
	qry := `select u.id, u.username from users u, tokens t where u.id = t.user_id and t.token = $1`
	row := a.DB.QueryRow(qry, refresh)
	if err := row.Scan(&subject, &username); err != nil {
		return handler.NewStatusError(http.StatusUnauthorized, "")
	}

	// Create access token
	access := createAccessToken(a, subject, username)

	// Response
	w.Write([]byte(access))
	return nil
}

func LogoutHandler(a *env.App, w http.ResponseWriter, r *http.Request) error {
	// Get refresh token from request
	refresh, err := getRefreshToken(r)
	if err != nil {
		return err
	}

	// Invalidate refresh token
	_, err = a.DB.Exec(`delete from tokens where token = $1`, refresh)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: check for access_token and blacklist in redis

	// Response
	return nil
}
