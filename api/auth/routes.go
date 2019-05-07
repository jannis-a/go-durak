package auth

import (
	"github.com/jannis-a/go-durak/routes"
)

var Routes = []routes.Route{
	{"login", "POST", "/login", LoginHandler},
	{"refresh", "GET", "/refresh", RefreshHandler},
	{"logout", "POST", "/logout", LogoutHandler},
}
