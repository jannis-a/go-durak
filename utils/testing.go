package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// API assert function for testing the response.
type ApiAssertFunc func(t *testing.T, req *http.Request, res *httptest.ResponseRecorder)

// API test case. Name is used for subtest. Data is the payload to send. Code is the
// status code of the response. Func the assert function for testing the response.
type ApiTestCase struct {
	Name string
	Data interface{}
	Code int
	Func ApiAssertFunc
}

func DispatchRequest(router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	req.RemoteAddr = "127.0.0.1"

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}
