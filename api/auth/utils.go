package auth

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/jannis-a/go-durak/env"
	"github.com/jannis-a/go-durak/handler"
)

const RefreshCookieName = "refresh_token"

func getRefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(RefreshCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", handler.NewStatusError(http.StatusUnauthorized, "")
		}
		return "", handler.NewStatusError(http.StatusBadRequest, "")
	}

	return cookie.Value, nil
}

func createAccessToken(a *env.App, user uint, username string) string {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.Itoa(int(user)),
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(a.Config.KEY))
	if err != nil {
		return ""
	}
	return signed
}

func ClaimsFromToken(a *env.App, r *http.Request) *Claims {
	header := r.Header.Get("Authorization")
	values := strings.Split(header, " ")
	if len(values) != 2 {
		return nil
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(values[1], claims, func(token *jwt.Token) (interface{}, error) {
		return a.Config.KEY, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if !token.Valid {
		return nil
	}
	return claims
}
