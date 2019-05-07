package auth

import (
	"github.com/jannis-a/go-durak/app"
)

var Routes = []app.Route{
	{"login", "POST", "/login", LoginHandler},
	{"refresh", "GET", "/refresh", RefreshHandler},
	{"logout", "POST", "/logout", LogoutHandler},
}
