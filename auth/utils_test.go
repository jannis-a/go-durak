package auth_test

import (
	"strconv"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"

	"github.com/jannis-a/go-durak/auth"
)

func TestCreateAccessToken(t *testing.T) {
	access := auth.CreateAccessToken(a, id, username)

	var claims auth.Claims
	token, err := jwt.ParseWithClaims(access, &claims, auth.KeyFunc(a))
	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.NotNil(t, claims)
	assert.Equal(t, strconv.Itoa(int(id)), claims.StandardClaims.Subject)
	assert.Equal(t, username, claims.Username)
}
