package brewerydb

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestHopList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/hop.list.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/hops/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	hl, err := client.Hop.List(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(hl.Hops) <= 0 {
		t.Fatal("Expected >0 hops")
	}

	c := "hop"
	for _, h := range hl.Hops {
		if c != h.Category {
			t.Fatalf("Hop Category = %s, wanted %s", h.Category, c)
		}
	}
}

func ExampleHopService_List() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	// Get a specific variety of hop with a given ID
	h, err := c.Hop.Get(84)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", h)

	// Get all types of hops
	hl, err := c.Hop.List(1)
	if err != nil {
		log.Fatal(err)
	}
	for _, h := range hl.Hops {
		fmt.Println(h.Name)
	}
}
