package auth

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/jannis-a/go-durak/app"
)

const RefreshCookieName = "refresh_token"

func getRefreshToken(r *http.Request) string {
	cookie, err := r.Cookie(RefreshCookieName)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func createAccessToken(a *app.App, user uint, username string) string {
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

func ClaimsFromToken(a *app.App, r *http.Request) *Claims {
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
