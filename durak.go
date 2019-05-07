package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"

	"github.com/jannis-a/go-durak/api/auth"
	"github.com/jannis-a/go-durak/api/users"
	"github.com/jannis-a/go-durak/env"
	"github.com/jannis-a/go-durak/routes"
)

func main() {
	// Create the application
	app := env.NewApp(nil)

	// Initialize routes
	routes.Register(app, "auth", auth.Routes)
	routes.Register(app, "users", users.Routes)

	// Display all available routes
	err := app.Router.Walk(routes.Walk)
	if err != nil {
		log.Fatal(err)
	}

	// Serve routes
	handler := handlers.LoggingHandler(os.Stdout, app.Router)
	log.Println("listening on", app.Config.BIND)
	log.Fatal(http.ListenAndServe(app.Config.BIND, handler))
}
