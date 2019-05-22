package users_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jannis-a/go-durak/users"
)

func TestValidatePassword(t *testing.T) {
	testsCases := []struct {
		name     string
		password string
		error    []string
	}{
		{"valid", "Val1dPa$$word", []string{}},
		{"invalid_len", "Inv4l!d", []string{users.ErrPasswordLen}},
		{"invalid_upper", "inv4lid!", []string{users.ErrPasswordUpper}},
		{"invalid_lower", "INV4LID!", []string{users.ErrPasswordLower}},
		{"invalid_number", "Invalid!", []string{users.ErrPasswordNumber}},
		{"invalid_special", "Inv4lidPassword", []string{users.ErrPasswordSpecial}},
		{"invalid_all", "", []string{
			users.ErrPasswordLen,
			users.ErrPasswordUpper,
			users.ErrPasswordLower,
			users.ErrPasswordNumber,
			users.ErrPasswordSpecial,
		}},
	}

	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, errors := users.ValidatePassword(tc.password)

			if len(tc.error) == 0 {
				assert.True(t, valid)
				assert.Nil(t, errors)
			} else {
				assert.False(t, valid)
				assert.Equal(t, tc.error, errors)
			}
		})
	}
}
