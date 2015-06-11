package brewerydb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

// Checks that the Request's URL query string contains name url-encoded with value=value.
func checkFormValue(t *testing.T, r *http.Request, name, value string) {
	if v := r.FormValue(name); v != value {
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

// Executes fn, expecting it to return an error
func testBadURL(t *testing.T, fn func() error) {
	origURL := apiURL
	apiURL = "http://%api.brewerydb.com/v2"
	if err := fn(); err == nil {
		t.Fatal("expected HTTP Request URL error")
	}
	apiURL = origURL
}

func TestNewRequest(t *testing.T) {
	setup()
	defer teardown()

	// `data` parameter should be a struct, not a string
	_, err := client.NewRequest("GET", "/heartbeat", "hello, world")
	if err == nil {
		t.Fatal("Expected query encoding error")
	}

	_, err = client.NewRequest("FOO", "/hearbeat", nil)
	if err == nil {
		t.Fatal("Expected HTTP method error")
	}
}

// for testing client.Do error handling
type testTransport struct{}

func (t testTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("fake round-trip error")
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	client.JSONWriter = ioutil.Discard

	const beerID = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{}")
	})

	_, err := client.Beer.Get(beerID)
	if err != nil {
		t.Fatal(err)
	}

	client.client.Transport = testTransport{}
	_, err = client.Beer.Get(beerID)
	if err == nil {
		t.Fatal("Expected net/http Do error")
	}
}

func TestYesNoUnmarshalJSON(t *testing.T) {
	q := struct {
		IsPrimary YesNo `url:"isPrimary"`
		IsClosed  YesNo `url:"isClosed,omitempty"`
	}{}

	js0 := []byte(`{"isPrimary":"Y"}`)

	if err := json.Unmarshal(js0, &q); err != nil {
		t.Error(err)
	}
	if v := YesNo(true); q.IsPrimary != v {
		t.Errorf("q.IsPrimary = %v, want %v", q.IsPrimary, v)
	}
	if v := YesNo(false); q.IsClosed != v {
		t.Errorf("q.IsClosed = %v, want %v", q.IsClosed, v)
	}

	js1 := []byte(`{"isPrimary":"N", "isClosed":"Y"}`)
	if err := json.Unmarshal(js1, &q); err != nil {
		t.Error(err)
	}
	if v := YesNo(false); q.IsPrimary != v {
		t.Errorf("q.IsPrimary = %v, want %v", q.IsPrimary, v)
	}
	if v := YesNo(true); q.IsClosed != v {
		t.Errorf("q.IsClosed = %v, want %v", q.IsClosed, v)
	}

	js2 := []byte(`{"isPrimary":"y", "isClosed":true}`)
	if err := json.Unmarshal(js2, &q); err == nil {
		t.Errorf(`Expected unmarshal error (only "Y" or "N" are valid YesNo JSON)`)
	}
}

// "What is Doppelbock?"
func Example_doppelbock() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	styles, err := c.Menu.Styles()
	if err != nil {
		panic(err)
	}
	for _, style := range styles {
		if style.ShortName == "Doppelbock" {
			fmt.Println("Doppelbock: \n", style.Description)
		}
	}
}

// "What is in Dragon's Milk?"
func Example_dragonsmilk() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	bl, err := c.Search.Beer("Dragon's Milk", nil)
	if err != nil {
		panic(err)
	}

	var beerID string
	for _, beer := range bl.Beers {
		if beer.Name == "Dragon's Milk" {
			beerID = beer.ID
		}
	}
	if beerID == "" {
		panic("Dragon's Milk not found")
	}

	ingredients, err := c.Beer.ListIngredients(beerID)
	if err != nil {
		panic(err)
	}

	adjuncts, err := c.Beer.ListAdjuncts(beerID)
	if err != nil {
		panic(err)
	}

	hops, err := c.Beer.ListHops(beerID)
	if err != nil {
		panic(err)
	}

	yeasts, err := c.Beer.ListYeasts(beerID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Dragon's Milk:")
	fmt.Println("  Ingredients:")
	for _, ingredient := range ingredients {
		fmt.Println("    " + ingredient.Name)
	}
	fmt.Println("  Adjuncts:")
	for _, adjunct := range adjuncts {
		fmt.Println("    " + adjunct.Name)
	}
	fmt.Println("  Hops:")
	for _, hop := range hops {
		fmt.Println("    " + hop.Name)
	}
	fmt.Println("  Yeasts:")
	for _, yeast := range yeasts {
		fmt.Println("    " + yeast.Name)
	}
}
