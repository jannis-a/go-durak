package app

import (
	"net/http"
)

// The Handler struct that takes a configured Env and a function matching our useful signature.
type Handler struct {
	app  *App
	Func HandlerFunc
}

// Define a type for handler functions with the App struct as parameter.
type HandlerFunc func(*App, http.ResponseWriter, *http.Request)

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Func(h.app, w, r)
}
