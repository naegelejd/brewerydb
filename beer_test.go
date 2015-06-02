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

	testBadURL(t, func() error {
		_, err := client.Beer.Get(id)
		return err
	})
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

	testBadURL(t, func() error {
		_, err := client.Beer.List(&BeerListRequest{ABV: abv})
		return err
	})
}

func makeTestBeer() *Beer {
	return &Beer{
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
		IsOrganic:       true,
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

}

func TestBeerAdd(t *testing.T) {
	setup()
	defer teardown()

	beer := makeTestBeer()

	const newID = "abcdef"
	mux.HandleFunc("/beers", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")

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
		checkPostFormValue(t, r, "isOrganic", "Y")
		checkPostFormValue(t, r, "label", beer.Label)
		checkPostFormValue(t, r, "brewery", beer.Brewery[0])
		checkPostFormValue(t, r, "servingTemperature", string(beer.ServingTemperature))
		checkPostFormValue(t, r, "availableId", strconv.Itoa(beer.AvailableID))
		checkPostFormValue(t, r, "srmId", strconv.Itoa(beer.SrmID))
		checkPostFormValue(t, r, "year", strconv.Itoa(beer.Year))

		// Check that fields tagged with "-" or "omitempty" are NOT encoded
		checkPostFormDNE(t, r, "id", "ID", "status", "Status", "beerVariationId", "BeerVariationID")

		fmt.Fprintf(w, `{"status":"...", "data":{"id":"%s"}, "message":"..."}`, newID)
	})

	id, err := client.Beer.Add(beer)
	if err != nil {
		t.Fatal(err)
	}
	if id != newID {
		t.Fatalf("new Beer ID = %v, want %v", id, newID)
	}

	_, err = client.Beer.Add(nil)
	if err == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		_, err = client.Beer.Add(beer)
		return err
	})
}

func TestBeerUpdate(t *testing.T) {
	setup()
	defer teardown()

	beer := makeTestBeer()
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		checkURLSuffix(t, r, beer.ID)

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
		checkPostFormValue(t, r, "isOrganic", "Y")
		checkPostFormValue(t, r, "label", beer.Label)
		checkPostFormValue(t, r, "brewery", beer.Brewery[0])
		checkPostFormValue(t, r, "servingTemperature", string(beer.ServingTemperature))
		checkPostFormValue(t, r, "availableId", strconv.Itoa(beer.AvailableID))
		checkPostFormValue(t, r, "srmId", strconv.Itoa(beer.SrmID))
		checkPostFormValue(t, r, "year", strconv.Itoa(beer.Year))

		// Check that fields tagged with "-" or "omitempty" are NOT encoded
		checkPostFormDNE(t, r, "id", "ID", "status", "Status", "beerVariationId", "BeerVariationID")
	})

	if err := client.Beer.Update(beer.ID, beer); err != nil {
		t.Fatal(err)
	}

	if client.Beer.Update(beer.ID, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Beer.Update(beer.ID, beer)
	})
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

	testBadURL(t, func() error {
		return client.Beer.Delete(id0)
	})
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

	testBadURL(t, func() error {
		return client.Beer.DeleteBrewery(beerID0, breweryID0, delReq)
	})
}

type beerDeleter func(string, int) error

func testBeerDeleteHelper(t *testing.T, name string, otherID int, del beerDeleter) {
	setup()
	defer teardown()

	const beerID = "o9TSOv"
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

	testBadURL(t, func() error {
		return del(beerID, otherID)
	})
}

func TestBeerDeleteAdjunct(t *testing.T) {
	adjunctID := 923
	testBeerDeleteHelper(t, "adjunct", adjunctID, client.Beer.DeleteAdjunct)
}

func TestBeerDeleteFermentable(t *testing.T) {
	fermentableID := 753
	testBeerDeleteHelper(t, "fermentable", fermentableID, client.Beer.DeleteFermentable)
}

func TestBeerDeleteHop(t *testing.T) {
	hopID := 42
	testBeerDeleteHelper(t, "hop", hopID, client.Beer.DeleteHop)
}

func TestBeerDeleteSocialAccount(t *testing.T) {
	socialID := 3
	testBeerDeleteHelper(t, "socialaccount", socialID, client.Beer.DeleteSocialAccount)
}

func TestBeerDeleteYeast(t *testing.T) {
	yeastID := 1835
	testBeerDeleteHelper(t, "yeast", yeastID, client.Beer.DeleteYeast)
}

func TestBeerListAdjuncts(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("adjunct.list.json", t)
	defer data.Close()

	const beerID = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "adjuncts" {
			t.Fatal("bad URL, expected \"/beer/:beerId/adjuncts\"")
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	al, err := client.Beer.ListAdjuncts(beerID)
	if err != nil {
		t.Fatal(err)
	}

	if len(al) <= 0 {
		t.Fatal("Expected >0 Adjuncts")
	}

	for _, a := range al {
		if a.ID <= 0 {
			t.Fatal("Expected ID >0")
		}
	}

	testBadURL(t, func() error {
		_, err := client.Beer.ListAdjuncts(beerID)
		return err
	})
}

func TestBeerListHops(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("hop.list.json", t)
	defer data.Close()

	const beerID = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "hops" {
			t.Fatal("bad URL, expected \"/beer/:beerId/hops\"")
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	al, err := client.Beer.ListHops(beerID)
	if err != nil {
		t.Fatal(err)
	}

	if len(al) <= 0 {
		t.Fatal("Expected >0 Hops")
	}

	for _, a := range al {
		if a.ID <= 0 {
			t.Fatal("Expected ID >0")
		}
	}

	testBadURL(t, func() error {
		_, err := client.Beer.ListHops(beerID)
		return err
	})
}

func TestBeerListIngredients(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("ingredient.list.json", t)
	defer data.Close()

	const beerID = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "ingredients" {
			t.Fatal("bad URL, expected \"/beer/:beerId/ingredients\"")
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	al, err := client.Beer.ListIngredients(beerID)
	if err != nil {
		t.Fatal(err)
	}

	if len(al) <= 0 {
		t.Fatal("Expected >0 Ingredients")
	}

	for _, a := range al {
		if a.ID <= 0 {
			t.Fatal("Expected ID >0")
		}
	}

	testBadURL(t, func() error {
		_, err := client.Beer.ListIngredients(beerID)
		return err
	})
}

func TestBeerListFermentables(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("fermentable.list.json", t)
	defer data.Close()

	const beerID = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "fermentables" {
			t.Fatal("bad URL, expected \"/beer/:beerId/fermentables\"")
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	al, err := client.Beer.ListFermentables(beerID)
	if err != nil {
		t.Fatal(err)
	}

	if len(al) <= 0 {
		t.Fatal("Expected >0 Fermentables")
	}

	for _, a := range al {
		if a.ID <= 0 {
			t.Fatal("Expected ID >0")
		}
	}

	testBadURL(t, func() error {
		_, err := client.Beer.ListFermentables(beerID)
		return err
	})
}

func TestBeerListVariations(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("beer.list.json", t)
	defer data.Close()

	const beerID = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "variations" {
			t.Fatal("bad URL, expected \"/beer/:beerId/variations\"")
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	al, err := client.Beer.ListVariations(beerID)
	if err != nil {
		t.Fatal(err)
	}

	if len(al) <= 0 {
		t.Fatal("Expected >0 Variations")
	}

	for _, a := range al {
		if l := 6; len(a.ID) != l {
			t.Fatalf("Variation ID len = %v, want %v", len(a.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Beer.ListVariations(beerID)
		return err
	})
}

func TestBeerListYeasts(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("yeast.list.json", t)
	defer data.Close()

	const beerID = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "yeasts" {
			t.Fatal("bad URL, expected \"/beer/:beerId/yeasts\"")
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	al, err := client.Beer.ListYeasts(beerID)
	if err != nil {
		t.Fatal(err)
	}

	if len(al) <= 0 {
		t.Fatal("Expected >0 Yeasts")
	}

	for _, a := range al {
		if a.ID <= 0 {
			t.Fatal("Expected ID >0")
		}
	}

	testBadURL(t, func() error {
		_, err := client.Beer.ListYeasts(beerID)
		return err
	})
}

type beerAdder func(string, int) error

func testBeerAddHelper(t *testing.T, name string, otherID int, add beerAdder) {
	setup()
	defer teardown()

	const beerID = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != name+"s" {
			t.Fatalf("bad URL, expected \"/beer/:beerId/%ss\"", name)
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, name+"Id", strconv.Itoa(otherID))
	})

	if err := add(beerID, otherID); err != nil {
		t.Fatal(err)
	}

	if add("******", otherID) == nil {
		t.Fatal("expected HTTP error")
	}

	testBadURL(t, func() error {
		return add(beerID, otherID)
	})
}

func TestBeerAddAdjunct(t *testing.T) {
	const adjunctID = 923
	testBeerAddHelper(t, "adjunct", adjunctID, client.Beer.AddAdjunct)
}

func TestBeerListBreweries(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("brewery.list.json", t)
	defer data.Close()

	const beerID = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "breweries" {
			t.Fatal("bad URL, expected \"/beer/:beerId/breweries\"")
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	bl, err := client.Beer.ListBreweries(beerID)
	if err != nil {
		t.Fatal(err)
	}

	if len(bl) <= 0 {
		t.Fatal("Expected >0 Breweries")
	}

	for _, b := range bl {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Brewery ID len = %d, wanted %d", len(b.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Beer.ListBreweries(beerID)
		return err
	})
}

func TestBeerAddBrewery(t *testing.T) {
	setup()
	defer teardown()

	const (
		beerID     = "o9TSOv"
		breweryID  = "jmGoBA"
		locationID = "z9H6HJ"
	)
	firstTest := true
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "breweries" {
			t.Fatal("bad URL, expected \"/beer/:beerId/breweries\"")
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "breweryId", breweryID)
		if firstTest {
			checkPostFormValue(t, r, "locationId", locationID)
		} else {
			checkPostFormDNE(t, r, "locationId")
		}
	})

	if err := client.Beer.AddBrewery(beerID, breweryID, &BeerBreweryRequest{locationID}); err != nil {
		t.Fatal(err)
	}

	firstTest = false
	if err := client.Beer.AddBrewery(beerID, breweryID, &BeerBreweryRequest{}); err != nil {
		t.Fatal(err)
	}

	if err := client.Beer.AddBrewery(beerID, breweryID, nil); err != nil {
		t.Fatal(err)
	}

	if client.Beer.AddBrewery("******", breweryID, nil) == nil {
		t.Fatal("expected HTTP 404 error")
	}

	testBadURL(t, func() error {
		return client.Beer.AddBrewery(beerID, breweryID, &BeerBreweryRequest{locationID})
	})
}

func TestBeerAddFermentable(t *testing.T) {
	const fermentableID = 753
	testBeerAddHelper(t, "fermentable", fermentableID, client.Beer.AddFermentable)
}

func TestBeerAddHop(t *testing.T) {
	const hopID = 42
	testBeerAddHelper(t, "hop", hopID, client.Beer.AddHop)
}

func TestBeerGetSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("beer.get.socialaccount.json", t)
	defer data.Close()

	const (
		beerID          = "o9TSOv"
		socialAccountID = 1
	)
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccount" {
			t.Fatal("bad URL, expected \"/beer/:beerId/socialaccount/:socialaccountId\"")
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(socialAccountID) {
			http.Error(w, "invalid SocialAccount ID", http.StatusNotFound)
		}
		io.Copy(w, data)

	})

	a, err := client.Beer.GetSocialAccount(beerID, socialAccountID)
	if err != nil {
		t.Fatal(err)
	}

	if a.ID != socialAccountID {
		t.Fatalf("SocialAccount ID = %v, want %v", a.ID, socialAccountID)
	}

	testBadURL(t, func() error {
		_, err := client.Beer.GetSocialAccount(beerID, socialAccountID)
		return err
	})
}

func TestBeerListSocialAccounts(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("beer.list.socialaccounts.json", t)
	defer data.Close()

	const beerID = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccounts" {
			t.Fatal("bad URL, expected \"/beer/:beerId/socialaccounts\"")
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	al, err := client.Beer.ListSocialAccounts(beerID)
	if err != nil {
		t.Fatal(err)
	}

	if len(al) <= 0 {
		t.Fatal("Expected >0 SocialAccounts")
	}

	for _, a := range al {
		if a.ID <= 0 {
			t.Fatal("Expected ID >0")
		}
	}

	testBadURL(t, func() error {
		_, err := client.Beer.ListSocialAccounts(beerID)
		return err
	})
}

func TestBeerAddSocialAccount(t *testing.T) {
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
		checkMethod(t, r, "POST")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccounts" {
			t.Fatal("bad URL, expected \"/beer/:beerId/socialaccounts\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "socialmediaId", strconv.Itoa(account.SocialMediaID))
		checkPostFormValue(t, r, "handle", account.Handle)

		checkPostFormDNE(t, r, "id", "ID", "socialMedia", "SocialSite")
	})

	if err := client.Beer.AddSocialAccount(id, account); err != nil {
		t.Fatal(err)
	}

	if client.Beer.AddSocialAccount("******", account) == nil {
		t.Fatal("expected HTTP error")
	}

	if client.Beer.AddSocialAccount(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Beer.AddSocialAccount(id, account)
	})
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
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}
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

	testBadURL(t, func() error {
		return client.Beer.UpdateSocialAccount(id, account)
	})
}

func TestBeerAddUPC(t *testing.T) {
	setup()
	defer teardown()

	const (
		beerID = "o9TSOv"
		upc    = 98765432100
	)
	fluidsizeID := 5
	firstTest := true
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "upcs" {
			t.Fatal("bad URL, expected \"/beer/:beerId/upcs\"")
		}
		if split[2] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "upcCode", fmt.Sprintf("%d", upc))
		if firstTest {
			checkPostFormValue(t, r, "fluidSizeId", strconv.Itoa(fluidsizeID))
		} else {
			checkPostFormDNE(t, r, "fluidSizeId")
		}

		// fluidsizeID is encoded as fluidSizeId... ensure other variants are nowhere to be found
		checkPostFormDNE(t, r, "upc", "Upc", "UPC", "fluidsizeId", "fluidsizeID")
	})

	if err := client.Beer.AddUPC(beerID, upc, &fluidsizeID); err != nil {
		t.Fatal(err)
	}

	firstTest = false
	if err := client.Beer.AddUPC(beerID, upc, nil); err != nil {
		t.Fatal(err)
	}

	if client.Beer.AddUPC("******", upc, nil) == nil {
		t.Fatal("expected HTTP 404 error")
	}

	testBadURL(t, func() error {
		return client.Beer.AddUPC(beerID, upc, &fluidsizeID)
	})
}

func TestBeerAddYeast(t *testing.T) {
	const yeastID = 1835
	testBeerAddHelper(t, "yeast", yeastID, client.Beer.AddYeast)
}

func TestBeerGetRandom(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/beer.get.random.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/beer/random", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")

		// TODO: check more request query values
		checkFormValue(t, r, "abv", "8")

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

	testBadURL(t, func() error {
		_, err := client.Beer.GetRandom(&RandomBeerRequest{ABV: "8"})
		return err
	})
}

func TestBeerListEvents(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("event.list.json", t)
	defer data.Close()

	const beerID = "o9TSOv"
	mux.HandleFunc("/beer/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "events" {
			t.Fatal("bad URL, expected \"/beer/:beerId/events\"")
		}
		if split[2] != beerID {
			http.Error(w, "invalid beer ID", http.StatusNotFound)
		}

		checkFormValue(t, r, "onlyWinners", "Y")

		io.Copy(w, data)

	})

	el, err := client.Beer.ListEvents(beerID, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(el) <= 0 {
		t.Fatal("Expected >0 Events")
	}

	for _, e := range el {
		if l := 6; l != len(e.ID) {
			t.Fatalf("Event ID len = %d, wanted %d", len(e.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Beer.ListEvents(beerID, false)
		return err
	})
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
