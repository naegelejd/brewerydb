package brewerydb

import (
	"net/http"
	"net/http/httptest"
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
