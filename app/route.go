package app

import (
	"fmt"

	"github.com/gorilla/mux"
)

type Route struct {
	Name    string
	Method  string
	Path    string
	Handler HandlerFunc
}

func (a *App) RegisterApi(prefix string, routes []Route) {
	router := a.Router.PathPrefix("/" + prefix).Subrouter()

	for _, r := range routes {
		router.
			Name(prefix + ":" + r.Name).
			Methods(r.Method).
			Path(r.Path).
			Handler(Handler{a, r.Handler})
	}
}

func Walk(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	methods, err := route.GetMethods()
	if err != nil {
		return nil
	}

	path, err := route.GetPathTemplate()
	if err != nil {
		return err
	}

	fmt.Println(methods, route.GetName(), path)
	return nil
}
