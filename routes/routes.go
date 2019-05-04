package routes

import (
	"fmt"

	"github.com/gorilla/mux"

	"github.com/jannis-a/go-durak/env"
	"github.com/jannis-a/go-durak/handler"
)

type Route struct {
	Name    string
	Method  string
	Path    string
	Handler handler.HandlerFunc
}

func Register(app *env.App, prefix string, routes []Route) {
	router := app.Router.PathPrefix("/" + prefix).Subrouter()

	for _, r := range routes {
		fn := handler.Handler{app, r.Handler}

		router.
			Name(prefix + ":" + r.Name).
			Methods(r.Method).
			Path(r.Path).
			Handler(fn)

	}
}

func Walk(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	methods, err := route.GetMethods()
	if err != nil {
		return nil
	}

	url, err := route.GetPathTemplate()
	if err != nil {
		return err
	}

	fmt.Println(methods, route.GetName(), url)
	return nil
}
