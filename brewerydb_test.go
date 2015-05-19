package brewerydb

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

var (
	mux     *http.ServeMux
	server  *httptest.Server
	client  *Client
	fakeKey = "abcdefghijklmnopqrstuvwxyz"
)

func loadTestData(filename string, t *testing.T) io.ReadCloser {
	data, err := os.Open("test_data/" + filename)
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	return data
}

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

// Checks that the HTTP Request's URL path ends with suffix, ignoring any trailing slashes.
func checkURLSuffix(t *testing.T, r *http.Request, suffix string) {
	if !strings.HasSuffix(strings.TrimSuffix(r.URL.Path, "/"), suffix) {
		t.Fatalf("URL path = %s, expected suffix = %s", r.URL.Path, suffix)
	}
}

// Checks that the Request's body contains name url-encoded with value=value
func checkPostFormValue(t *testing.T, r *http.Request, name, value string) {
	if v := r.PostFormValue(name); v != value {
		t.Fatalf("%s = %v, want %v", name, v, value)
	}
}

// Checks that each key is NOT url-encoded in the Request's Body
func checkPostFormDNE(t *testing.T, r *http.Request, keys ...string) {
	if err := r.ParseForm(); err != nil {
		t.Fatal(err)
	}
	formMap := map[string][]string(r.PostForm)
	for _, key := range keys {
		if _, ok := formMap[key]; ok {
			t.Fatalf("form value '%s' should not be encoded", key)
		}
	}
}
