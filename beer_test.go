package brewerydb

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestBeerGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("beer.get.json", t)
	defer data.Close()

	const id = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, id)
		io.Copy(w, data)
	})

	b, err := client.Beer.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if b.ID != id {
		t.Fatalf("Beer ID = %v, want %v", b.ID, id)
	}
}

func TestBeerList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("beer.list.json", t)
	defer data.Close()

	const abv = "8"
	mux.HandleFunc("/beers/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		if v := r.FormValue("abv"); v != abv {
			t.Fatalf("Request.FormValue abv = %v, wanted %v", v, abv)
		}
		// TODO: check more request query values
		io.Copy(w, data)
	})

	bl, err := client.Beer.List(&BeerListRequest{ABV: abv})
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Beers) <= 0 {
		t.Fatal("Expected >0 beers")
	}
	for _, b := range bl.Beers {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Beer ID len = %d, wanted %d", len(b.ID), l)
		}
	}
}

func TestBeerAdd(t *testing.T) {

}

func TestBeerUpdate(t *testing.T) {

}

func TestBeerDelete(t *testing.T) {
	setup()
	defer teardown()

	const (
		id0 = "o9TSOv"
		id1 = "******"
	)
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "beer" {
			t.Fatal("bad URL, expected \"/beer/:beerId\"")
		}
		if split[2] != id0 {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}
		// TODO: should Delete tests care about JSON response?
		// io.Copy(w, bytes.NewBufferString(`{
		// "status": "success",
		// "data": 1,
		// "message": "..."
		// }`))

	})

	if err := client.Beer.Delete(id0); err != nil {
		t.Fatal(err)
	}

	if err := client.Beer.Delete(id1); err == nil {
		t.Fatal("expected HTTP 404")
	}
}

func TestBeerDeleteBrewery(t *testing.T) {
	setup()
	defer teardown()

	const (
		beerID0    = "o9TSOv"
		beerID1    = "******"
		breweryID0 = "jmGoBA"
		breweryID1 = "~~~~~~"
		locationID = "z9H6HJ"
	)
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "beer" || split[3] != "brewery" {
			t.Fatal("bad URL, expected \"/beer/:beerId/brewery/:breweryId\"")
		}
		if split[2] != beerID0 {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}
		if split[4] != breweryID0 {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}

		if v := r.FormValue("locationid"); v != "" && v != locationID {
			t.Fatalf("Request.FormValue locationId = %v, wanted %v", v, locationID)
		}
	})

	delReq := &BeerBreweryRequest{LocationID: locationID}

	// check valid DeleteBrewery with locationID
	if err := client.Beer.DeleteBrewery(beerID0, breweryID0, delReq); err != nil {
		t.Fatal(err)
	}
	// check valid DeleteBrewery with nil locationID
	if err := client.Beer.DeleteBrewery(beerID0, breweryID0, nil); err != nil {
		t.Fatal(err)
	}

	if client.Beer.DeleteBrewery(beerID0, breweryID1, nil) == nil {
		t.Fatal("expected HTTP 404 error")
	}

	if client.Beer.DeleteBrewery(beerID1, breweryID0, nil) == nil {
		t.Fatal("expected HTTP 404 error")
	}
}

type beerDeleter func(string, int) error

func testBeerDeleteHelper(t *testing.T, name, beerID string, otherID int, del beerDeleter) {
	setup()
	defer teardown()

	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "beer" || split[3] != name {
			t.Fatalf("bad URL, expected \"/beer/:beerId/%s/:%sId\"", name, name)
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(otherID) {
			http.Error(w, "invalid "+name+" ID", http.StatusNotFound)
		}
	})

	if err := del(beerID, otherID); err != nil {
		log.Fatal(err)
	}

	if del(beerID, -1) == nil {
		t.Fatal("expected HTTP 404 error")
	}

	if del("*~*~*~", otherID) == nil {
		t.Fatal("expected HTTP 404 error")
	}
}

func TestBeerDeleteAdjunct(t *testing.T) {
	setup()
	defer teardown()

	beerID := "o9TSOv"
	adjunctID := 923
	testBeerDeleteHelper(t, "adjunct", beerID, adjunctID, client.Beer.DeleteAdjunct)
}

func TestBeerDeleteFermentable(t *testing.T) {
	setup()
	defer teardown()

	beerID := "o9TSOv"
	fermentableID := 753
	testBeerDeleteHelper(t, "fermentable", beerID, fermentableID, client.Beer.DeleteFermentable)
}

func TestBeerDeleteHop(t *testing.T) {
	setup()
	defer teardown()

	beerID := "o9TSOv"
	hopID := 42
	testBeerDeleteHelper(t, "hop", beerID, hopID, client.Beer.DeleteHop)
}

func TestBeerDeleteSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	beerID := "o9TSOv"
	socialID := 3
	testBeerDeleteHelper(t, "socialaccount", beerID, socialID, client.Beer.DeleteSocialAccount)
}

func TestBeerDeleteYeast(t *testing.T) {
	setup()
	defer teardown()

	beerID := "o9TSOv"
	yeastID := 1835
	testBeerDeleteHelper(t, "yeast", beerID, yeastID, client.Beer.DeleteYeast)
}

func TestBeerRandom(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/beer.getrandom.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/beer/random/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	b, err := client.Beer.GetRandom(&RandomBeerRequest{ABV: "8"})
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

func ExampleBeerService_List() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	// Get first 40 beers with an ABV between 8.0 and 9.0, descending, alphabetical
	bl, err := c.Beer.List(&BeerListRequest{ABV: "8", Sort: SortDescending})
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range bl.Beers {
		fmt.Println(b.Name, b.ID)
	}
}

func ExampleBeerService_Breweries() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	// Get breweries for a given beer
	breweries, err := c.Beer.ListBreweries("jmGoBA")
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range breweries {
		fmt.Println(b.Name)
	}
}

// Get a random beer with an ABV between 8.0 and 9.0
func ExampleBeerService_Random() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	req := &RandomBeerRequest{
		ABV: "8",
	}
	b, err := c.Beer.GetRandom(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b.Name)
	fmt.Println(b.Style.Name)
	fmt.Println(b.Labels.Large)
}
