package user

import (
	"github.com/go-chi/chi"

	"github.com/jannis-a/go-durak/app"
)

type App struct {
	*app.App
}

func Routes(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", ListHandler(app))
	r.Post("/", CreateHandler(app))
	r.Get("/{username}", DetailHandler(app))
	r.Patch("/{username}", UpdateHandler(app))
	r.Delete("/{username}", DeleteHandler(app))

	return r
}
