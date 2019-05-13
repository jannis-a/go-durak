package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"

	"github.com/jannis-a/go-durak/app"
	"github.com/jannis-a/go-durak/auth"
	"github.com/jannis-a/go-durak/users"
)

func main() {
	// Create app and register routes
	a := app.NewApp()
	a.RegisterApi("auth", auth.Routes)
	a.RegisterApi("users", users.Routes)

	// Display bind address and all available routes
	fmt.Println("Listening on", a.Config.BIND)
	a.WalkRoutes()

	// Serve routes
	handler := handlers.LoggingHandler(os.Stdout, a.Router)
	log.Fatal(http.ListenAndServe(a.Config.BIND, handler))
}
