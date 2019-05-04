package utils

import (
	"encoding/json"
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
	if vars := mux.Vars(r); vars != nil {
		return vars[name]
	}
	return ""
}

func RenderJson(w http.ResponseWriter, value interface{}) {
	encoded, err := json.Marshal(&value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(encoded)
}
