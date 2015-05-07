package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestBeerRandom(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/beer.random.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/beer/random/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		io.Copy(w, data)
	})

	b, err := client.Beer.Random(&RandomBeerRequest{ABV: "8"})
	if err != nil {
		t.Fatal(err)
	}

	// Can't really verify specific information since it's a random beer
	if len(b.Name) <= 0 {
		t.Fatal("Expected non-empty beer name")
	}
	if len(b.ID) <= 0 {
		t.Fatal("Expected non-empty beer ID")
	}
}
