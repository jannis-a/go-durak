package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jannis-a/go-durak/database"
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

	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.Recoverer)

	addr := fmt.Sprintf("%s:%d", viper.GetString("HOST"), viper.GetInt("PORT"))
	log.Printf("listening on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
