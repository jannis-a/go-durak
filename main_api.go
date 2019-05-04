package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jannis-a/go-durak/api/users"
	"github.com/jannis-a/go-durak/env"
	"github.com/jannis-a/go-durak/routes"
)

func main() {
	// Create the application
	app := env.NewApp(nil)

	// Initialize routes
	routes.Register(app, "users", users.Routes)

	// Display all available routes
	err := app.Router.Walk(routes.Walk)
	if err != nil {
		log.Fatal(err)
	}

	// Serve routes
	addr := ":" + strconv.Itoa(app.Config.PORT)
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
