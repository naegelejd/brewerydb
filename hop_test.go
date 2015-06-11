package brewerydb

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestHopGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("hop.get.json", t)
	defer data.Close()

	const id = 42
	mux.HandleFunc("/hop/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, strconv.Itoa(id))
		io.Copy(w, data)
	})

	h, err := client.Hop.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if h.ID != id {
		t.Fatalf("Hop ID = %v, want %v", h.ID, id)
	}

	testBadURL(t, func() error {
		_, err := client.Hop.Get(id)
		return err
	})
}

func TestHopList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("hop.list.json", t)
	defer data.Close()

	const page = 1
	mux.HandleFunc("/hops", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		io.Copy(w, data)
	})

	hl, err := client.Hop.List(page)
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

	testBadURL(t, func() error {
		_, err := client.Hop.List(page)
		return err
	})
}

// Get a specific variety of hop with a given ID
func ExampleHopService_Get() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	h, err := c.Hop.Get(84)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", h)
}

// Get all types of hops
func ExampleHopService_List() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	hl, err := c.Hop.List(1)
	if err != nil {
		panic(err)
	}
	for _, h := range hl.Hops {
		fmt.Println(h.Name)
	}
}
