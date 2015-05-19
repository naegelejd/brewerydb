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
	mux.HandleFunc("/beers", func(w http.ResponseWriter, r *http.Request) {
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
	setup()
	defer teardown()

	beer := &Beer{
		ID:              "o9TSOv",
		Name:            "The Truth",
		Description:     "Hop bomb",
		FoodPairings:    "Barbecue",
		OriginalGravity: "1.0",
		ABV:             "8.7",
		IBU:             "80",
		GlasswareID:     5,
		Glass:           Glass{ID: 5, Name: "Pint"},
		StyleID:         31,
		IsOrganic:       "N",
		Labels: Images{
			"https://s3.amazonaws.com/brewerydbapi/beer/o9TSOv/upload_nIhalb-icon.png",
			"https://s3.amazonaws.com/brewerydbapi/beer/o9TSOv/upload_nIhalb-medium.png",
			"https://s3.amazonaws.com/brewerydbapi/beer/o9TSOv/upload_nIhalb-large.png",
		},
		Label:              "https://s3.amazonaws.com/brewerydbapi/beer/o9TSOv/upload_nIhalb-large.png",
		Brewery:            []string{"jmGoBA"},
		ServingTemperature: TemperatureCool,
		Status:             "verified",
		AvailableID:        1,
		Available:          Availability{ID: 1, Name: "Year Round"},
		SrmID:              6,
		SRM:                SRM{ID: 6, Name: "6", Hex: "F8A600"},
		Year:               2013,
	}

	const id = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		checkURLSuffix(t, r, id)

		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}

		checkPostFormValue(t, r, "name", beer.Name)
		checkPostFormValue(t, r, "description", beer.Description)
		checkPostFormValue(t, r, "foodPairings", beer.FoodPairings)
		checkPostFormValue(t, r, "originalGravity", beer.OriginalGravity)
		checkPostFormValue(t, r, "abv", beer.ABV)
		checkPostFormValue(t, r, "ibu", beer.IBU)
		checkPostFormValue(t, r, "glasswareId", strconv.Itoa(beer.GlasswareID))
		checkPostFormValue(t, r, "styleId", strconv.Itoa(beer.StyleID))
		checkPostFormValue(t, r, "isOrganic", beer.IsOrganic)
		checkPostFormValue(t, r, "label", beer.Label)
		checkPostFormValue(t, r, "brewery", beer.Brewery[0])
		checkPostFormValue(t, r, "servingTemperature", string(beer.ServingTemperature))
		checkPostFormValue(t, r, "availableId", strconv.Itoa(beer.AvailableID))
		checkPostFormValue(t, r, "srmId", strconv.Itoa(beer.SrmID))
		checkPostFormValue(t, r, "year", strconv.Itoa(beer.Year))

		// Check that fields tagged with "-" or "omitempty" are NOT encoded
		checkPostFormDNE(t, r, "id", "ID", "status", "Status", "beerVariationId", "BeerVariationID")
	})

	if err := client.Beer.Update(id, beer); err != nil {
		t.Fatal(err)
	}

	if client.Beer.Update(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}
}

func TestBeerUpdateSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	account := &SocialAccount{
		ID:            2,
		SocialMediaID: 8,
		SocialSite: SocialSite{
			ID:      8,
			Name:    "Google Plus",
			Website: "https://plus.google.com/u/0/",
		},
		Handle: "flying_dog",
	}

	const id = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccount" {
			t.Fatal("bad URL, expected \"/beer/:beerId/socialaccount/:socialaccountId\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(account.ID) {
			http.Error(w, "invalid SocialAccount ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "socialmediaId", strconv.Itoa(account.SocialMediaID))
		checkPostFormValue(t, r, "handle", account.Handle)

		checkPostFormDNE(t, r, "id", "socialMedia", "SocialSite")
	})

	if err := client.Beer.UpdateSocialAccount(id, account); err != nil {
		t.Fatal(err)
	}

	if client.Beer.UpdateSocialAccount("******", account) == nil {
		t.Fatal("expected HTTP error")
	}

	if client.Beer.UpdateSocialAccount(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}
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
