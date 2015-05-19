package brewerydb

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestEventGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("event.get.json", t)
	defer data.Close()

	const id = "0oZVAo"
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, id)
		io.Copy(w, data)
	})

	e, err := client.Event.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if e.ID != id {
		t.Fatalf("Event ID = %v, want %v", e.ID, id)
	}
}

func TestEventList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("event.list.json", t)
	defer data.Close()

	const year = 2015
	mux.HandleFunc("/events/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		if v := r.FormValue("year"); v != strconv.Itoa(year) {
			t.Fatalf("Request.FormValue year = %v, wanted %v", v, year)
		}
		// TODO: check more request query values
		io.Copy(w, data)
	})

	el, err := client.Event.List(&EventListRequest{Year: year})
	if err != nil {
		t.Fatal(err)
	}
	if len(el.Events) <= 0 {
		t.Fatal("Expected >0 events")
	}
	for _, e := range el.Events {
		if l := 6; l != len(e.ID) {
			t.Fatalf("Event ID len = %d, wanted %d", len(e.ID), l)
		}
	}
}

func TestEventDelete(t *testing.T) {
	setup()
	defer teardown()

	const id = "0oZVAo"
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "event" {
			t.Fatal("bad URL, expected \"/event/:eventId\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
	})

	if err := client.Event.Delete(id); err != nil {
		t.Fatal(err)
	}

	if err := client.Event.Delete("******"); err == nil {
		t.Fatal("expected HTTP 404")
	}
}

func TestEventDeleteBeer(t *testing.T) {
	setup()
	defer teardown()

	const (
		eventID = "0oZVAo"
		beerID  = "o9TSOv"
	)
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "event" || split[3] != "beer" {
			t.Fatal("bad URL, expected \"/event/:eventId/beer/:beerId\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}
	})

	if err := client.Event.DeleteBeer(eventID, beerID); err != nil {
		t.Fatal(err)
	}

	if err := client.Event.DeleteBeer("******", beerID); err == nil {
		t.Fatal("expected HTTP 404")
	}

	if err := client.Event.DeleteBeer(eventID, "~~~~~~"); err == nil {
		t.Fatal("expected HTTP 404")
	}
}

func TestEventDeleteBrewery(t *testing.T) {
	setup()
	defer teardown()

	const (
		eventID   = "0oZVAo"
		breweryID = "jmGoBA"
	)
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "event" || split[3] != "brewery" {
			t.Fatal("bad URL, expected \"/event/:eventId/brewery/:breweryId\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}
	})

	if err := client.Event.DeleteBrewery(eventID, breweryID); err != nil {
		t.Fatal(err)
	}

	if err := client.Event.DeleteBrewery("******", breweryID); err == nil {
		t.Fatal("expected HTTP 404")
	}

	if err := client.Event.DeleteBrewery(eventID, "~~~~~~"); err == nil {
		t.Fatal("expected HTTP 404")
	}
}

func TestEventDeleteAwardCategory(t *testing.T) {
	setup()
	defer teardown()

	const (
		eventID         = "0oZVAo"
		awardCategoryID = 2
	)
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "event" || split[3] != "awardcategory" {
			t.Fatal("bad URL, expected \"/event/:eventId/awardcategory/:awardcategoryId\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(awardCategoryID) {
			http.Error(w, "invalid awardcategory ID", http.StatusNotFound)
		}
	})

	if err := client.Event.DeleteAwardCategory(eventID, awardCategoryID); err != nil {
		t.Fatal(err)
	}

	if err := client.Event.DeleteAwardCategory("******", awardCategoryID); err == nil {
		t.Fatal("expected HTTP 404")
	}

	if err := client.Event.DeleteAwardCategory(eventID, -1); err == nil {
		t.Fatal("expected HTTP 404")
	}
}

func TestEventDeleteAwardPlace(t *testing.T) {
	setup()
	defer teardown()

	const (
		eventID      = "0oZVAo"
		awardPlaceID = 2
	)
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "event" || split[3] != "awardplace" {
			t.Fatal("bad URL, expected \"/event/:eventId/awardplace/:awardplaceId\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(awardPlaceID) {
			http.Error(w, "invalid awardplace ID", http.StatusNotFound)
		}
	})

	if err := client.Event.DeleteAwardPlace(eventID, awardPlaceID); err != nil {
		t.Fatal(err)
	}

	if err := client.Event.DeleteAwardPlace("******", awardPlaceID); err == nil {
		t.Fatal("expected HTTP 404")
	}

	if err := client.Event.DeleteAwardPlace(eventID, -1); err == nil {
		t.Fatal("expected HTTP 404")
	}
}

func TestEventDeleteAwardSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	const (
		eventID  = "0oZVAo"
		socialID = 2
	)
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "event" || split[3] != "socialaccount" {
			t.Fatal("bad URL, expected \"/event/:eventId/socialaccount/:socialaccountId\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(socialID) {
			http.Error(w, "invalid socialaccount ID", http.StatusNotFound)
		}
	})

	if err := client.Event.DeleteSocialAccount(eventID, socialID); err != nil {
		t.Fatal(err)
	}

	if err := client.Event.DeleteSocialAccount("******", socialID); err == nil {
		t.Fatal("expected HTTP 404")
	}

	if err := client.Event.DeleteSocialAccount(eventID, -1); err == nil {
		t.Fatal("expected HTTP 404")
	}
}
