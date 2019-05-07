package users

import (
	"github.com/jannis-a/go-durak/app"
)

var Routes = []app.Route{
	{"create", "POST", "", CreateHandler},
	{"list", "GET", "", ListHandler},
	{"read", "GET", "/{username}", DetailHandler},
	{"update", "PATCH", "/{username}", UpdateHandler},
	{"delete", "DELETE", "/{username}", DeleteHandler},
}
