package users

import (
	"github.com/jannis-a/go-durak/app"
)

var Routes = []app.Route{
	{Name: "create", Method: "POST", Handler: CreateHandler},
	{Name: "list", Method: "GET", Handler: ListHandler},
	{Name: "read", Method: "GET", Path: "/{username}", Handler: DetailHandler},
	{Name: "update", Method: "PATCH", Path: "/{username}", Handler: UpdateHandler},
	{Name: "delete", Method: "DELETE", Path: "/{username}", Handler: DeleteHandler},
}
