package brewerydb

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestSearchBeer(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("search.beer.json", t)
	defer data.Close()

	const (
		query = "flying"
		page  = 1
	)
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		checkFormValue(t, r, "q", query)
		// TODO: check more request query values
		io.Copy(w, data)
	})

	bl, err := client.Search.Beer(query, &SearchRequest{Page: page})
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Beers) <= 0 {
		t.Fatal("Expected >0 beers")
	}
	for _, b := range bl.Beers {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Beer ID len = %d, want %d", len(b.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Search.Beer(query, &SearchRequest{Page: page})
		return err
	})
}

func TestSearchBrewery(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("search.brewery.json", t)
	defer data.Close()

	const (
		query = "dog"
		page  = 1
	)
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		checkFormValue(t, r, "q", query)
		// TODO: check more request query values
		io.Copy(w, data)
	})

	bl, err := client.Search.Brewery(query, &SearchRequest{Page: page})
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Breweries) <= 0 {
		t.Fatal("Expected >0 breweries")
	}
	for _, b := range bl.Breweries {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Brewery ID len = %d, want %d", len(b.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Search.Brewery(query, &SearchRequest{Page: page})
		return err
	})
}

func TestSearchEvent(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("search.event.json", t)
	defer data.Close()

	const (
		query = "festival"
		page  = 1
	)
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)

		checkFormValue(t, r, "q", query)
		// TODO: check more request query values
		io.Copy(w, data)
	})

	bl, err := client.Search.Event(query, &SearchRequest{Page: page})
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Events) <= 0 {
		t.Fatal("Expected >0 events")
	}
	for _, b := range bl.Events {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Event ID len = %d, want %d", len(b.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Search.Event(query, &SearchRequest{Page: page})
		return err
	})
}

func TestSearchGuild(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("search.guild.json", t)
	defer data.Close()

	const (
		query = "maryland"
		page  = 1
	)
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)

		checkFormValue(t, r, "q", query)
		// TODO: check more request query values

		io.Copy(w, data)
	})

	bl, err := client.Search.Guild(query, &SearchRequest{Page: page})
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Guilds) <= 0 {
		t.Fatal("Expected >0 guilds")
	}
	for _, b := range bl.Guilds {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Guild ID len = %d, want %d", len(b.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Search.Guild(query, &SearchRequest{Page: page})
		return err
	})
}

func TestSearchStyle(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("search.style.json", t)
	defer data.Close()

	const (
		query = "Pale Ale"
	)
	mux.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, "style")

		checkFormValue(t, r, "q", "Pale Ale")
		checkFormValue(t, r, "withDescriptions", "Y")

		io.Copy(w, data)
	})

	sl, err := client.Search.Style(query, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(sl) <= 0 {
		t.Fatal("Expected >0 styles")
	}
	for _, s := range sl {
		if s.ID <= 0 {
			t.Fatal("Expected ID >0")
		}
	}

	testBadURL(t, func() error {
		_, err := client.Search.Style(query, true)
		return err
	})
}

func TestSearchGeoPoint(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("search.geopoint.json", t)
	defer data.Close()

	const (
		latitude  = 35.772096
		longitude = -78.638614
	)
	mux.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[2] != "geo" || split[3] != "point" {
			t.Fatal("bad URL, expected \"/search/geo/point\"")
		}

		checkFormValue(t, r, "lat", fmt.Sprintf("%f", latitude))
		checkFormValue(t, r, "lng", fmt.Sprintf("%f", longitude))
		// TODO: check more form values

		io.Copy(w, data)
	})

	req := &GeoPointRequest{Latitude: latitude, Longitude: longitude}
	ll, err := client.Search.GeoPoint(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(ll) <= 0 {
		t.Fatal("Expected >0 Locations")
	}
	for _, l := range ll {
		if len(l.ID) != 6 {
			t.Fatal("Expected ID len to be 6")
		}
	}

	testBadURL(t, func() error {
		_, err := client.Search.GeoPoint(req)
		return err
	})
}

func TestSearchUPC(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("search.upc.json", t)
	defer data.Close()

	const (
		code = 606905008303
	)
	firstTest := true
	mux.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, "upc")

		checkFormValue(t, r, "code", fmt.Sprintf("%d", code))

		if firstTest {
			io.Copy(w, data)
		} else {
			fmt.Fprint(w, `{"status":"success", "data":[]}`)
		}
	})

	bl, err := client.Search.UPC(code)
	if err != nil {
		t.Fatal(err)
	}
	if len(bl) <= 0 {
		t.Fatal("Expected >0 Beers")
	}
	for _, b := range bl {
		if len(b.ID) != 6 {
			t.Fatal("Expected ID len to be 6")
		}
	}

	firstTest = false
	bl, _ = client.Search.UPC(code)
	if bl != nil {
		t.Fatal("Expected nil []Beer")
	}

	testBadURL(t, func() error {
		_, err := client.Search.UPC(code)
		return err
	})
}

func TestMakeActualSearchRequest(t *testing.T) {
	const (
		p  = 0
		q  = "good beer"
		tp = searchBeer
	)
	req := makeActualSearchRequest(nil, q, tp)
	if req.Page != p {
		t.Fatalf("Page = %v, want %v", req.Page, p)
	}
	if req.Query != q {
		t.Fatalf("Query = %v, want %v", req.Query, q)
	}
	if req.Type != tp {
		t.Fatalf("Type = %v, want %v", req.Type, tp)
	}
}
