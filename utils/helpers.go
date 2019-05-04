package utils

import (
	"go/build"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

func GetPackagePath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	return path.Join(gopath, "src", "github.com", "jannis-a", "go-durak")
}

func GetRouteParam(r *http.Request, name string) string {
	vars := mux.Vars(r)
	return vars[name]
}
