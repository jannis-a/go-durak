package user

import (
	"github.com/go-chi/chi"

	"github.com/jannis-a/go-durak/config"
)

func Routes(c *config.Config) *chi.Mux {
	c.Db.AutoMigrate(&User{})

	r := chi.NewRouter()

	r.Get("/", ListHandler(c))
	r.Post("/", CreateHandler(c))
	r.Get("/{username}", DetailHandler(c))
	r.Patch("/{username}", UpdateHandler(c))
	r.Delete("/{username}", DeleteHandler(c))

	return r
}
