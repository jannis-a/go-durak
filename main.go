package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/jannis-a/go-durak/config"
	"github.com/jannis-a/go-durak/user"
)

func Routes(c *config.Config) *chi.Mux {
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
	configuration, err := config.New()
	if err != nil {
		log.Panicln("Configuration error", err)
	}

	defer configuration.Db.Close()

	router := Routes(configuration)

	log.Println("listening on:", configuration.Constants.PORT)
	log.Fatal(http.ListenAndServe(":"+configuration.Constants.PORT, router))
}
