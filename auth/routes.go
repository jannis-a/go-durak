package auth

import (
	"github.com/jannis-a/go-durak/app"
)

var Routes = []app.Route{
	{Name: "login", Method: "POST", Path: "/login", Handler: LoginHandler},
	{Name: "refresh", Method: "GET", Path: "/refresh", Handler: RefreshHandler},
	{Name: "logout", Method: "POST", Path: "/logout", Handler: LogoutHandler},
}
