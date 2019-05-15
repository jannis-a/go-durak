package app

import (
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

type Route struct {
	Name    string
	Method  string
	Path    string
	Handler HandlerFunc
	Query   []string
}

func (a *App) RegisterApi(prefix string, routes []Route) {
	router := a.Router.PathPrefix("/" + prefix).Subrouter()

	for _, r := range routes {
		route := router.
			Name(prefix + ":" + r.Name).
			Methods(r.Method).
			Path(r.Path).
			Handler(Handler{a, r.Handler})

		if len(r.Query) > 0 {
			route.Queries(r.Query...)
		}
	}
}

func (a *App) WalkRoutes() {
	routes := tablewriter.NewWriter(os.Stdout)
	routes.SetHeader([]string{"methods", "name", "path"})
	routes.SetAlignment(tablewriter.ALIGN_LEFT)
	routes.SetBorder(false)

	walkFunc := func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}

		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		routes.Append([]string{strings.Join(methods, ","), route.GetName(), path})
		return nil
	}

	err := a.Router.Walk(walkFunc)
	if err != nil {
		log.Error(err)
	} else {
		routes.Render()
	}
}
