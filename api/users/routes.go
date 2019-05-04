package users

import (
	"github.com/jannis-a/go-durak/routes"
)

var Routes = []routes.Route{
	{"create", "POST", "", CreateHandler},
	{"list", "GET", "", ListHandler},
	{"read", "GET", "/{username}", DetailHandler},
	{"update", "PATCH", "/{username}", UpdateHandler},
	{"delete", "DELETE", "/{username}", DeleteHandler},
}
