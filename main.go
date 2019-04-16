package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/jannis-a/go-durak/app"
	"github.com/jannis-a/go-durak/user"
)

func Routes(c *app.App) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/users", user.Routes(c))
	})

	return router
}

func main() {
	application, err := app.New()
	if err != nil {
		log.Panicln("Configuration error", err)
	}

	router := Routes(application)

	log.Println("listening on:", application.Config.PORT)
	log.Fatal(http.ListenAndServe(":"+application.Config.PORT, router))
}
