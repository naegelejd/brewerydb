package brewerydb

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var (
	mux     *http.ServeMux
	server  *httptest.Server
	client  *Client
	fakeKey = "abcdefghijklmnopqrstuvwxyz"
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	apiURL = server.URL
	client = NewClient(fakeKey)
}

func teardown() {
	server.Close()
}

// Checks that the HTTP Request's method matches the given method.
func checkMethod(t *testing.T, r *http.Request, method string) {
	if method != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, method)
	}
}

// Checks that the HTTP Request contains a key "p" with a value matching the given page.
func checkPage(t *testing.T, r *http.Request, page int) {
	if p := r.FormValue("p"); p != strconv.Itoa(page) {
		t.Fatalf("Request.FormValue p = %v, want %v", p, page)
	}
}
