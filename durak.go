package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"github.com/jannis-a/go-durak/app"
	"github.com/jannis-a/go-durak/auth"
	"github.com/jannis-a/go-durak/users"
)

func main() {
	// Create app and register routes
	a := app.NewApp()
	a.Register("auth", auth.Routes)
	a.Register("users", users.Routes)

	// Display all available routes
	err := a.Router.Walk(app.Walk)
	if err != nil {
		log.Fatal(err)
	}

	// Serve routes
	handler := handlers.LoggingHandler(os.Stdout, a.Router)
	log.Println("listening on", a.Config.BIND)
	log.Fatal(http.ListenAndServe(a.Config.BIND, handler))
}
