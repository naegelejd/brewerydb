package brewerydb

import (
	"io"
	"net/http"
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
		if q := r.FormValue("q"); q != query {
			t.Fatalf("Request.FormValue q = %v, want %v", q, query)
		}
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
		if q := r.FormValue("q"); q != query {
			t.Fatalf("Request.FormValue q = %v, want %v", q, query)
		}
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
		if q := r.FormValue("q"); q != query {
			t.Fatalf("Request.FormValue q = %v, want %v", q, query)
		}
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
		if q := r.FormValue("q"); q != query {
			t.Fatalf("Request.FormValue q = %v, want %v", q, query)
		}
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
