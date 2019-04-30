package users

import (
	"github.com/go-chi/chi"

	"github.com/jannis-a/go-durak/app"
)

func Routes(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/", CreateHandler(app))

	r.Get("/", ListHandler(app))
	r.Get("/{username}", DetailHandler(app))

	r.Get("/users/me", DetailHandler(app))
	r.Patch("/users/me", UpdateHandler(app))
	r.Delete("/users/me", DeleteHandler(app))

	return r
}
