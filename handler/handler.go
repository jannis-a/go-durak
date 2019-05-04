package handler

import (
	"log"
	"net/http"

	"github.com/jannis-a/go-durak/env"
)

// Define a type for handler functions with the App struct as parameter.
type HandlerFunc func(*env.App, http.ResponseWriter, *http.Request) error

// The Handler struct that takes a configured Env and a function matching our useful signature.
type Handler struct {
	*env.App
	Func HandlerFunc
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.Func(h.App, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific HTTP status code.
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default to serving a HTTP 500
			http.Error(w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}
	}
}
