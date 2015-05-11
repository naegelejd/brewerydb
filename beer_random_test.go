package brewerydb

import (
	"fmt"
	"io"
	"log"
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
		checkMethod(t, r, "GET")
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

// Get a random beer with an ABV between 8.0 and 9.0
func ExampleBeerService_Random() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	req := &RandomBeerRequest{
		ABV: "8",
	}
	b, err := c.Beer.Random(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b.Name)
	fmt.Println(b.Style.Name)
	fmt.Println(b.Labels.Large)
}
