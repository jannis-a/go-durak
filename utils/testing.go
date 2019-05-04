package utils

import (
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

func DispatchRequest(router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}
