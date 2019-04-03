package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jannis-a/go-durak/database"
	"github.com/jannis-a/go-durak/handlers"
	"github.com/jannis-a/go-durak/models"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config: %s\n", err))
	}

	db := database.Open()
	defer db.Close()
	db.AutoMigrate(&models.User{})

	app := handlers.NewApp(db)

	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recoverer)

	router.Route("/users", func(r chi.Router) {
		r.Get("/", app.UserList)
		r.Post("/", app.UserCreate)

		r.Route("/{username}", func(r chi.Router) {
			r.Get("/", app.UserDetail)
			r.Patch("/", app.UserUpdate)
			r.Delete("/", app.UserDelete)
		})
	})
	addr := fmt.Sprintf("%s:%d", viper.GetString("HOST"), viper.GetInt("PORT"))
	log.Printf("listening on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
