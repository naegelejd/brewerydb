package brewerydb

import (
	"fmt"
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

	testBadURL(t, func() error {
		_, err := client.Event.Get(id)
		return err
	})
}

func TestEventList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("event.list.json", t)
	defer data.Close()

	const year = 2015
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
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

	testBadURL(t, func() error {
		_, err := client.Event.List(&EventListRequest{Year: year})
		return err
	})
}

func makeTestEvent() *Event {
	return &Event{
		ID:             "k2jMtH",
		Name:           "Yellowstone Beer Fest",
		Type:           EventFestival,
		StartDate:      "2015-07-18",
		EndDate:        "2015-07-18",
		Description:    "Regional beer fest in Cody, Wyoming",
		Year:           "2015",
		Time:           "from 3:00 PM - 8:00 PM",
		Price:          "$30 - $35",
		VenueName:      "Park County Complex",
		StreetAddress:  "1501 Stampede Ave.",
		Locality:       "Cody",
		Region:         "Wyoming",
		PostalCode:     "82414",
		CountryISOCode: "US",
		Website:        "http://www.yellowstonebeerfest.com/",
		Longitude:      -109.05873,
		Latitude:       44.520076,
		Image:          "https://s3.amazonaws.com/brewerydbapi/event/0oZVAo/upload_KjVkrq-large.png",
		Images: Images{
			"https://s3.amazonaws.com/brewerydbapi/event/0oZVAo/upload_KjVkrq-icon.png",
			"https://s3.amazonaws.com/brewerydbapi/event/0oZVAo/upload_KjVkrq-medium.png",
			"https://s3.amazonaws.com/brewerydbapi/event/0oZVAo/upload_KjVkrq-large.png",
		},
	}
}

func TestEventAdd(t *testing.T) {
	setup()
	defer teardown()

	event := makeTestEvent()

	const newID = "abcdef"
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")

		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}

		checkPostFormValue(t, r, "name", event.Name)
		checkPostFormValue(t, r, "type", string(EventFestival))
		checkPostFormValue(t, r, "startDate", event.StartDate)
		checkPostFormValue(t, r, "endDate", event.EndDate)
		checkPostFormValue(t, r, "description", event.Description)
		checkPostFormValue(t, r, "time", event.Time)
		checkPostFormValue(t, r, "time", event.Time)
		checkPostFormValue(t, r, "price", event.Price)
		checkPostFormValue(t, r, "venueName", event.VenueName)
		checkPostFormValue(t, r, "streetAddress", event.StreetAddress)
		checkPostFormValue(t, r, "locality", event.Locality)
		checkPostFormValue(t, r, "region", event.Region)
		checkPostFormValue(t, r, "postalCode", event.PostalCode)
		checkPostFormValue(t, r, "countryIsoCode", event.CountryISOCode)
		checkPostFormValue(t, r, "website", event.Website)
		checkPostFormValue(t, r, "longitude", fmt.Sprintf("%.5f", event.Longitude))
		checkPostFormValue(t, r, "latitude", fmt.Sprintf("%f", event.Latitude))
		checkPostFormValue(t, r, "image", event.Image)

		// Check that fields tagged with "-" or "omitempty" are NOT encoded
		checkPostFormDNE(t, r, "id", "ID", "images", "Images", "status",
			"Status", "extendedAddress", "ExtendedAddress", "phone", "Phone")

		fmt.Fprintf(w, `{"status":"...", "data":{"id":"%s"}, "message":"..."}`, newID)
	})

	id, err := client.Event.Add(event)
	if err != nil {
		t.Fatal(err)
	}
	if id != newID {
		t.Fatalf("new Event ID = %v, want %v", id, newID)
	}

	_, err = client.Event.Add(nil)
	if err == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		_, err := client.Event.Add(event)
		return err
	})
}

func TestEventUpdate(t *testing.T) {
	setup()
	defer teardown()

	event := makeTestEvent()
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		checkURLSuffix(t, r, event.ID)

		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}

		checkPostFormValue(t, r, "name", event.Name)
		checkPostFormValue(t, r, "type", string(EventFestival))
		checkPostFormValue(t, r, "startDate", event.StartDate)
		checkPostFormValue(t, r, "endDate", event.EndDate)
		checkPostFormValue(t, r, "description", event.Description)
		checkPostFormValue(t, r, "time", event.Time)
		checkPostFormValue(t, r, "time", event.Time)
		checkPostFormValue(t, r, "price", event.Price)
		checkPostFormValue(t, r, "venueName", event.VenueName)
		checkPostFormValue(t, r, "streetAddress", event.StreetAddress)
		checkPostFormValue(t, r, "locality", event.Locality)
		checkPostFormValue(t, r, "region", event.Region)
		checkPostFormValue(t, r, "postalCode", event.PostalCode)
		checkPostFormValue(t, r, "countryIsoCode", event.CountryISOCode)
		checkPostFormValue(t, r, "website", event.Website)
		checkPostFormValue(t, r, "longitude", fmt.Sprintf("%.5f", event.Longitude))
		checkPostFormValue(t, r, "latitude", fmt.Sprintf("%f", event.Latitude))
		checkPostFormValue(t, r, "image", event.Image)

		// Check that fields tagged with "-" or "omitempty" are NOT encoded
		checkPostFormDNE(t, r, "id", "ID", "images", "Images", "status",
			"Status", "extendedAddress", "ExtendedAddress", "phone", "Phone")
	})

	if err := client.Event.Update(event.ID, event); err != nil {
		t.Fatal(err)
	}

	if client.Event.Update(event.ID, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Event.Update(event.ID, event)
	})
}

func TestEventGetAwardCategory(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("event.get.awardcategory.json", t)
	defer data.Close()

	const (
		eventID           = "cJio9R"
		awardCategoryID   = 87
		awardCategoryName = "American IPA"
	)
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "awardcategory" {
			t.Fatal("bad URL, expected \"/event/:eventId/awardcategory/:awardcategoryId\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(awardCategoryID) {
			http.Error(w, "invalid AwardCategory ID", http.StatusNotFound)
		}
		io.Copy(w, data)

	})

	a, err := client.Event.GetAwardCategory(eventID, awardCategoryID)
	if err != nil {
		t.Fatal(err)
	}

	if a.ID != awardCategoryID {
		t.Fatalf("AwardCategory ID = %v, want %v", a.ID, awardCategoryID)
	}
	if a.Name != awardCategoryName {
		t.Fatalf("AwardCategory Name = %v, want %v", a.Name, awardCategoryName)
	}

	testBadURL(t, func() error {
		_, err := client.Event.GetAwardCategory(eventID, awardCategoryID)
		return err
	})
}

func TestEventListAwardCategory(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("event.list.awardcategories.json", t)
	defer data.Close()

	const eventID = "cJio9R"
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "awardcategories" {
			t.Fatal("bad URL, expected \"/event/:eventId/awardcategories\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	al, err := client.Event.ListAwardCategories(eventID)
	if err != nil {
		t.Fatal(err)
	}

	if len(al) <= 0 {
		t.Fatal("Expected >0 AwardCategories")
	}

	for _, a := range al {
		if len(a.Name) <= 0 {
			t.Fatal("Expected non-empty AwardCategory Name")
		}
	}

	testBadURL(t, func() error {
		_, err := client.Event.ListAwardCategories(eventID)
		return err
	})
}

func makeTestAwardCategory() *AwardCategory {
	return &AwardCategory{
		ID:          1,
		Name:        "Best in Show",
		Description: "Best Brew",
		Image:       "http://www.fakeimage.com/1.jpg",
	}
}

func TestEventAddAwardCategory(t *testing.T) {
	setup()
	defer teardown()

	category := makeTestAwardCategory()

	const id = "k2jMtH"
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "awardcategories" {
			t.Fatal("bad URL, expected \"/event/:eventId/awardcategories\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "name", category.Name)
		checkPostFormValue(t, r, "description", category.Description)
		checkPostFormValue(t, r, "image", category.Image)

		checkPostFormDNE(t, r, "id", "ID", "CreateDate", "UpdateDate")
	})

	if err := client.Event.AddAwardCategory(id, category); err != nil {
		t.Fatal(err)
	}

	if client.Event.AddAwardCategory(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Event.AddAwardCategory(id, category)
	})
}

func TestEventUpdateAwardCategory(t *testing.T) {
	setup()
	defer teardown()

	category := makeTestAwardCategory()

	const id = "k2jMtH"
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "awardcategory" {
			t.Fatal("bad URL, expected \"/event/:eventId/awardcategory/:awardcategoryId\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(category.ID) {
			http.Error(w, "invalid AwardCategory ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "name", category.Name)
		checkPostFormValue(t, r, "description", category.Description)
		checkPostFormValue(t, r, "image", category.Image)

		checkPostFormDNE(t, r, "id", "ID", "CreateDate", "UpdateDate")
	})

	if err := client.Event.UpdateAwardCategory(id, category); err != nil {
		t.Fatal(err)
	}

	if client.Event.UpdateAwardCategory(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Event.UpdateAwardCategory(id, category)
	})
}

func TestEventGetAwardPlace(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("event.get.awardplace.json", t)
	defer data.Close()

	const (
		eventID        = "cJio9R"
		awardPlaceID   = 3
		awardPlaceName = "Silver"
	)
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "awardplace" {
			t.Fatal("bad URL, expected \"/event/:eventId/awardplace/:awardplaceId\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(awardPlaceID) {
			http.Error(w, "invalid AwardPlace ID", http.StatusNotFound)
		}
		io.Copy(w, data)

	})

	a, err := client.Event.GetAwardPlace(eventID, awardPlaceID)
	if err != nil {
		t.Fatal(err)
	}

	if a.ID != awardPlaceID {
		t.Fatalf("AwardPlace ID = %v, want %v", a.ID, awardPlaceID)
	}
	if a.Name != awardPlaceName {
		t.Fatalf("AwardPlace Name = %v, want %v", a.Name, awardPlaceName)
	}

	testBadURL(t, func() error {
		_, err := client.Event.GetAwardPlace(eventID, awardPlaceID)
		return err
	})
}

func TestEventListAwardPlace(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("event.list.awardplaces.json", t)
	defer data.Close()

	const eventID = "cJio9R"
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "awardplaces" {
			t.Fatal("bad URL, expected \"/event/:eventId/awardplaces\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	al, err := client.Event.ListAwardPlaces(eventID)
	if err != nil {
		t.Fatal(err)
	}

	if len(al) <= 0 {
		t.Fatal("Expected >0 AwardPlaces")
	}

	for _, a := range al {
		if len(a.Name) <= 0 {
			t.Fatal("Expected non-empty AwardPlace Name")
		}
	}

	testBadURL(t, func() error {
		_, err := client.Event.ListAwardPlaces(eventID)
		return err
	})
}

func makeTestAwardPlace() *AwardPlace {
	return &AwardPlace{
		ID:          1,
		Name:        "First",
		Description: "First Place",
		Image:       "http://www.fakeimage.com/2.jpg",
	}
}

func TestEventAddAwardPlace(t *testing.T) {
	setup()
	defer teardown()

	place := makeTestAwardPlace()
	const id = "k2jMtH"
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "awardplaces" {
			t.Fatal("bad URL, expected \"/event/:eventId/awardplaces\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "name", place.Name)
		checkPostFormValue(t, r, "description", place.Description)
		checkPostFormValue(t, r, "image", place.Image)

		checkPostFormDNE(t, r, "id", "ID", "CreateDate", "UpdateDate")
	})

	if err := client.Event.AddAwardPlace(id, place); err != nil {
		t.Fatal(err)
	}

	if client.Event.AddAwardPlace(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Event.AddAwardPlace(id, place)
	})
}

func TestEventUpdateAwardPlace(t *testing.T) {
	setup()
	defer teardown()

	place := makeTestAwardPlace()
	const id = "k2jMtH"
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "awardplace" {
			t.Fatal("bad URL, expected \"/event/:eventId/awardplace/:awardplaceId\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(place.ID) {
			http.Error(w, "invalid AwardPlace ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "name", place.Name)
		checkPostFormValue(t, r, "description", place.Description)
		checkPostFormValue(t, r, "image", place.Image)

		checkPostFormDNE(t, r, "id", "ID", "CreateDate", "UpdateDate")
	})

	if err := client.Event.UpdateAwardPlace(id, place); err != nil {
		t.Fatal(err)
	}

	if client.Event.UpdateAwardPlace(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Event.UpdateAwardPlace(id, place)
	})
}

func TestEventGetBeer(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("beer.get.json", t)
	defer data.Close()

	const (
		eventID = "k2jMtH"
		beerID  = "o9TSOv"
	)
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "beer" {
			t.Fatal("bad URL, expected \"/event/:eventId/beer/:beerId\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}
		io.Copy(w, data)

	})

	b, err := client.Event.GetBeer(eventID, beerID)
	if err != nil {
		t.Fatal(err)
	}

	if b.ID != beerID {
		t.Fatalf("Beer ID = %v, want %v", b.ID, beerID)
	}

	testBadURL(t, func() error {
		_, err := client.Event.GetBeer(eventID, beerID)
		return err
	})
}

func TestEventListBeer(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("beer.list.json", t)
	defer data.Close()

	const (
		eventID = "k2jMtH"
		page    = 3
	)
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "beers" {
			t.Fatal("bad URL, expected \"/event/:eventId/beers\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}

		checkFormValue(t, r, "onlyWinners", "Y")
		checkFormValue(t, r, "awardcategoryId", "2")
		checkFormValue(t, r, "awardplaceId", "3")

		io.Copy(w, data)

	})

	req := &EventBeersRequest{Page: page, OnlyWinners: "Y", AwardPlaceID: 3, AwardCategoryID: 2}
	bl, err := client.Event.ListBeers(eventID, req)
	if err != nil {
		t.Fatal(err)
	}

	if len(bl.Beers) <= 0 {
		t.Fatal("Expected >0 Beers")
	}

	for _, b := range bl.Beers {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Beer ID len = %d, wanted %d", len(b.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Event.ListBeers(eventID, req)
		return err
	})
}

func makeTestEventChangeBeerRequest() *EventChangeBeerRequest {
	return &EventChangeBeerRequest{
		IsPouring:       "Y",
		AwardCategoryID: 2,
		AwardPlaceID:    3,
	}

}

func TestEventAddBeer(t *testing.T) {
	setup()
	defer teardown()

	change := makeTestEventChangeBeerRequest()
	const (
		eventID = "k2jMtH"
		beerID  = "o9TSOv"
	)
	firstTest := true
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "beers" {
			t.Fatal("bad URL, expected \"/event/:eventId/beers\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}

		if firstTest {
			checkPostFormValue(t, r, "isPouring", change.IsPouring)
			checkPostFormValue(t, r, "awardcategoryId", strconv.Itoa(change.AwardCategoryID))
			checkPostFormValue(t, r, "awardplaceId", strconv.Itoa(change.AwardPlaceID))
		}
	})

	if err := client.Event.AddBeer(eventID, beerID, change); err != nil {
		t.Fatal(err)
	}

	// Allowed to pass nil *EventChangeBeerRequest
	firstTest = false
	if err := client.Event.AddBeer(eventID, beerID, nil); err != nil {
		t.Fatal(err)
	}

	testBadURL(t, func() error {
		return client.Event.AddBeer(eventID, beerID, change)
	})
}

func TestEventUpdateBeer(t *testing.T) {
	setup()
	defer teardown()

	change := makeTestEventChangeBeerRequest()
	const (
		eventID = "k2jMtH"
		beerID  = "o9TSOv"
	)
	firstTest := true
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "beer" {
			t.Fatal("bad URL, expected \"/event/:eventId/beer/:beerId\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != beerID {
			http.Error(w, "invalid Beer ID", http.StatusNotFound)
		}

		if firstTest {
			checkPostFormValue(t, r, "isPouring", change.IsPouring)
			checkPostFormValue(t, r, "awardcategoryId", strconv.Itoa(change.AwardCategoryID))
			checkPostFormValue(t, r, "awardplaceId", strconv.Itoa(change.AwardPlaceID))
		}
	})

	if err := client.Event.UpdateBeer(eventID, beerID, change); err != nil {
		t.Fatal(err)
	}

	// Allowed to pass nil *EventChangeBeerRequest
	firstTest = false
	if err := client.Event.UpdateBeer(eventID, beerID, nil); err != nil {
		t.Fatal(err)
	}

	testBadURL(t, func() error {
		return client.Event.UpdateBeer(eventID, beerID, change)
	})
}

func TestEventGetBrewery(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("brewery.get.json", t)
	defer data.Close()

	const (
		eventID   = "k2jMtH"
		breweryID = "jmGoBA"
	)
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "brewery" {
			t.Fatal("bad URL, expected \"/event/:eventId/brewery/:breweryId\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}
		io.Copy(w, data)

	})

	b, err := client.Event.GetBrewery(eventID, breweryID)
	if err != nil {
		t.Fatal(err)
	}

	if b.ID != breweryID {
		t.Fatalf("Brewery ID = %v, want %v", b.ID, breweryID)
	}

	testBadURL(t, func() error {
		_, err := client.Event.GetBrewery(eventID, breweryID)
		return err
	})
}

func TestEventListBreweries(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("brewery.list.json", t)
	defer data.Close()

	const (
		eventID = "k2jMtH"
		page    = 3
	)
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "breweries" {
			t.Fatal("bad URL, expected \"/event/:eventId/breweries\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}

		checkFormValue(t, r, "onlyWinners", "Y")
		checkFormValue(t, r, "awardcategoryId", "3")
		checkFormValue(t, r, "awardplaceId", "2")

		io.Copy(w, data)
	})

	req := &EventBreweriesRequest{Page: page, OnlyWinners: "Y", AwardCategoryID: 3, AwardPlaceID: 2}
	bl, err := client.Event.ListBreweries(eventID, req)
	if err != nil {
		t.Fatal(err)
	}

	if len(bl.Breweries) <= 0 {
		t.Fatal("Expected >0 Breweries")
	}

	for _, b := range bl.Breweries {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Brewery ID len = %d, wanted %d", len(b.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Event.ListBreweries(eventID, req)
		return err
	})
}

func makeTestEventChangeBreweryRequest() *EventChangeBreweryRequest {
	return &EventChangeBreweryRequest{
		AwardCategoryID: 2,
		AwardPlaceID:    3,
	}
}

func TestEventAddBrewery(t *testing.T) {
	setup()
	defer teardown()

	change := makeTestEventChangeBreweryRequest()
	const (
		eventID   = "k2jMtH"
		breweryID = "jmGoBA"
	)
	firstTest := true
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "breweries" {
			t.Fatal("bad URL, expected \"/event/:eventId/breweries\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}

		if firstTest {
			checkPostFormValue(t, r, "awardcategoryId", strconv.Itoa(change.AwardCategoryID))
			checkPostFormValue(t, r, "awardplaceId", strconv.Itoa(change.AwardPlaceID))
		}
	})

	if err := client.Event.AddBrewery(eventID, breweryID, change); err != nil {
		t.Fatal(err)
	}

	// Allowed to pass nil *EventChangeBreweryRequest
	firstTest = false
	if err := client.Event.AddBrewery(eventID, breweryID, nil); err != nil {
		t.Fatal(err)
	}

	testBadURL(t, func() error {
		return client.Event.AddBrewery(eventID, breweryID, change)
	})
}

func TestEventUpdateBrewery(t *testing.T) {
	setup()
	defer teardown()

	change := makeTestEventChangeBreweryRequest()
	const (
		eventID   = "k2jMtH"
		breweryID = "jmGoBA"
	)
	firstTest := true
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "brewery" {
			t.Fatal("bad URL, expected \"/event/:eventId/brewery/:breweryId\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}

		if firstTest {
			checkPostFormValue(t, r, "awardcategoryId", strconv.Itoa(change.AwardCategoryID))
			checkPostFormValue(t, r, "awardplaceId", strconv.Itoa(change.AwardPlaceID))
		}
	})

	if err := client.Event.UpdateBrewery(eventID, breweryID, change); err != nil {
		t.Fatal(err)
	}

	// Allowed to pass nil *EventChangeBreweryRequest
	firstTest = false
	if err := client.Event.UpdateBrewery(eventID, breweryID, nil); err != nil {
		t.Fatal(err)
	}

	testBadURL(t, func() error {
		return client.Event.UpdateBrewery(eventID, breweryID, change)
	})
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

	testBadURL(t, func() error {
		return client.Event.Delete(id)
	})
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

	testBadURL(t, func() error {
		return client.Event.DeleteBeer(eventID, beerID)
	})
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

	testBadURL(t, func() error {
		return client.Event.DeleteBrewery(eventID, breweryID)
	})
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

	testBadURL(t, func() error {
		return client.Event.DeleteAwardCategory(eventID, awardCategoryID)
	})
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

	testBadURL(t, func() error {
		return client.Event.DeleteAwardPlace(eventID, awardPlaceID)
	})
}

func TestEventGetSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	// TODO: acquire socialaccounts for Events
	data := loadTestData("beer.get.socialaccount.json", t)
	defer data.Close()

	const (
		eventID         = "k2jMtH"
		socialAccountID = 1
	)
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccount" {
			t.Fatal("bad URL, expected \"/event/:eventId/socialaccount/:socialaccountId\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(socialAccountID) {
			http.Error(w, "invalid SocialAccount ID", http.StatusNotFound)
		}
		io.Copy(w, data)

	})

	a, err := client.Event.GetSocialAccount(eventID, socialAccountID)
	if err != nil {
		t.Fatal(err)
	}

	if a.ID != socialAccountID {
		t.Fatalf("SocialAccount ID = %v, want %v", a.ID, socialAccountID)
	}

	testBadURL(t, func() error {
		_, err := client.Event.GetSocialAccount(eventID, socialAccountID)
		return err
	})
}

func TestEventListSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	// TODO: acquire socialaccounts for Events
	data := loadTestData("beer.list.socialaccounts.json", t)
	defer data.Close()

	const eventID = "k2jMtH"
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccounts" {
			t.Fatal("bad URL, expected \"/event/:eventId/socialaccounts\"")
		}
		if split[2] != eventID {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	al, err := client.Event.ListSocialAccounts(eventID)
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
		_, err := client.Event.ListSocialAccounts(eventID)
		return err
	})
}

func TestEventAddSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	account := &SocialAccount{
		ID:            2,
		SocialMediaID: 4,
		SocialSite: SocialSite{
			ID:      4,
			Name:    "Untappd",
			Website: "https://www.untappd.com",
		},
		Handle: "yellowstone_beer_fest",
	}

	const id = "k2jMtH"
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccounts" {
			t.Fatal("bad URL, expected \"/event/:eventId/socialaccounts\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "socialmediaId", strconv.Itoa(account.SocialMediaID))
		checkPostFormValue(t, r, "handle", account.Handle)

		checkPostFormDNE(t, r, "id", "ID", "socialMedia", "SocialSite")
	})

	if err := client.Event.AddSocialAccount(id, account); err != nil {
		t.Fatal(err)
	}

	if client.Event.AddSocialAccount("******", account) == nil {
		t.Fatal("expected HTTP error")
	}

	if client.Event.AddSocialAccount(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Event.AddSocialAccount(id, account)
	})
}

func TestEventUpdateSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	account := &SocialAccount{
		ID:            2,
		SocialMediaID: 4,
		SocialSite: SocialSite{
			ID:      4,
			Name:    "Untappd",
			Website: "https://www.untappd.com",
		},
		Handle: "yellowstone_beer_fest",
	}

	const id = "k2jMtH"
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccount" {
			t.Fatal("bad URL, expected \"/event/:eventId/socialaccount/:socialaccountId\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Event ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(account.ID) {
			http.Error(w, "invalid SocialAccount ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "socialmediaId", strconv.Itoa(account.SocialMediaID))
		checkPostFormValue(t, r, "handle", account.Handle)

		checkPostFormDNE(t, r, "id", "ID", "socialMedia", "SocialSite")
	})

	if err := client.Event.UpdateSocialAccount(id, account); err != nil {
		t.Fatal(err)
	}

	if client.Event.UpdateSocialAccount("******", account) == nil {
		t.Fatal("expected HTTP error")
	}

	if client.Event.UpdateSocialAccount(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Event.UpdateSocialAccount(id, account)
	})
}

func TestEventDeleteSocialAccount(t *testing.T) {
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

	testBadURL(t, func() error {
		return client.Event.DeleteSocialAccount(eventID, socialID)
	})
}
