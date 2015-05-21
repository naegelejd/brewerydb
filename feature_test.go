package brewerydb

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestFeatureGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("feature.get.json", t)
	defer data.Close()

	mux.HandleFunc("/featured/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	f, err := client.Feature.Get()
	if err != nil {
		t.Fatal(err)
	}
	if id := 117; f.ID != id {
		t.Fatalf("Feature ID = %v, want %v", f.ID, id)
	}
	if breweryID := "DXvlfF"; f.BreweryID != breweryID {
		t.Fatalf("Feature BreweryID = %v, want %v", f.BreweryID, breweryID)
	}
	if beerID := "79SG11"; f.BeerID != beerID {
		t.Fatalf("Feature BeerID = %v, want %v", f.BeerID, beerID)
	}

	testBadURL(t, func() error {
		_, err := client.Feature.Get()
		return err
	})
}

func TestFeatureList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("feature.list.json", t)
	defer data.Close()

	const (
		page = 1
		year = 2015
	)
	mux.HandleFunc("/features", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		if v := r.FormValue("year"); v != strconv.Itoa(year) {
			t.Fatalf("Request.FormValue year = %v, wanted %v", v, year)
		}
		// TODO: check more request query values
		io.Copy(w, data)
	})

	fl, err := client.Feature.List(&FeatureListRequest{Page: page, Year: year})
	if err != nil {
		t.Fatal(err)
	}
	if len(fl.Features) <= 0 {
		t.Fatal("Expected >0 features")
	}

	for _, f := range fl.Features {
		if l := 6; l != len(f.Beer.ID) {
			t.Fatalf("Features Beer.ID len = %d, wanted %d", len(f.Beer.ID), l)
		}
		if l := 6; l != len(f.Brewery.ID) {
			t.Fatalf("Features Brewery.ID len = %d, wanted %d", len(f.Brewery.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Feature.List(nil)
		return err
	})
}

func TestFeatureByWeek(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("feature.byweek.json", t)
	defer data.Close()

	const (
		year = 2015
		week = 7
	)
	mux.HandleFunc("/feature/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, fmt.Sprintf("%4d-%02d", year, week))
		io.Copy(w, data)
	})

	f, err := client.Feature.ByWeek(year, week)
	if err != nil {
		t.Fatal(err)
	}
	if id := 101; f.ID != id {
		t.Fatalf("Feature ID = %v, want %v", f.ID, id)
	}
	if breweryID := "ouUm8w"; f.BreweryID != breweryID {
		t.Fatalf("Feature BreweryID = %v, want %v", f.BreweryID, breweryID)
	}
	if beerID := "fA0O8C"; f.BeerID != beerID {
		t.Fatalf("Feature BeerID = %v, want %v", f.BeerID, beerID)
	}

	testBadURL(t, func() error {
		_, err := client.Feature.ByWeek(year, week)
		return err
	})
}
