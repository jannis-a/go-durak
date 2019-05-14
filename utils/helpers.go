package utils

import (
	"crypto/rand"
	"encoding/json"
	"go/build"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gorilla/mux"
)

func GetPackagePath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	return path.Join(gopath, "src", "github.com", "jannis-a", "go-durak")
}

func HttpError(w http.ResponseWriter, code int, text string) {
	if text == "" {
		text = http.StatusText(code)
	}

	http.Error(w, text, code)
}

func GetIpAddr(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func GetRouteParam(r *http.Request, name string) string {
	vars := mux.Vars(r)
	if vars == nil {
		return ""
	}
	return vars[name]
}

func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func RenderErrors(w http.ResponseWriter, ae url.Values) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	RenderJson(w, map[string]url.Values{"errors": ae})
}

func RenderJson(w http.ResponseWriter, value interface{}) {
	data, err := json.Marshal(&value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}
