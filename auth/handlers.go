package auth

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/raja/argon2pw"

	"github.com/jannis-a/go-durak/app"
	"github.com/jannis-a/go-durak/utils"
)

func LoginHandler(a *app.App, w http.ResponseWriter, r *http.Request) {
	var (
		creds    Credentials
		subject  uint
		password string
	)

	// Decode credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		utils.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if creds.Username == "" || creds.Password == "" {
		utils.HttpError(w, http.StatusBadRequest, "")
		return
	}

	// Get user data and
	row := a.DB.QueryRow(`select id, password from users where username = $1`, creds.Username)
	if err := row.Scan(&subject, &password); err != nil {
		utils.HttpError(w, http.StatusUnauthorized, "")
		return
	}

	// Validate credentials
	valid, err := argon2pw.CompareHashWithPassword(password, creds.Password)
	if err != nil || !valid {
		utils.HttpError(w, http.StatusUnauthorized, "")
		return
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
	access := CreateAccessToken(a, subject, creds.Username)

	// Response
	http.SetCookie(w, &http.Cookie{
		Name:  RefreshCookieName,
		Value: refresh,
		// TODO: Secure:   true,
		HttpOnly: true,
	})
	w.Write([]byte(access))
}

func RefreshHandler(a *app.App, w http.ResponseWriter, r *http.Request) {
	var (
		subject  uint
		username string
	)

	// Get refresh token from request
	refresh := getRefreshToken(r)
	if refresh == "" {
		utils.HttpError(w, http.StatusUnauthorized, "")
		return
	}

	// Validate token and fetch subject data
	qry := `select u.id, u.username from users u, tokens t where u.id = t.user_id and t.token = $1`
	row := a.DB.QueryRow(qry, refresh)
	if err := row.Scan(&subject, &username); err != nil {
		utils.HttpError(w, http.StatusUnauthorized, "")
		return
	}

	// Create access token
	access := CreateAccessToken(a, subject, username)

	// Response
	w.Write([]byte(access))
}

func LogoutHandler(a *app.App, w http.ResponseWriter, r *http.Request) {
	// Get refresh token from request
	refresh := getRefreshToken(r)
	if refresh == "" {
		utils.HttpError(w, http.StatusUnauthorized, "")
		return
	}

	// Invalidate refresh token
	_, err := a.DB.Exec(`delete from tokens where token = $1`, refresh)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: check for access_token and blacklist in redis

	// Response
	return
}
