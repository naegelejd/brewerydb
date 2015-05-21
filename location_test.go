package brewerydb

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestLocationGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("location.get.json", t)
	defer data.Close()

	const id = "z9H6HJ"
	mux.HandleFunc("/location/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, id)
		io.Copy(w, data)
	})

	l, err := client.Location.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if l.ID != id {
		t.Fatalf("Location ID = %v, want %v", l.ID, id)
	}

	testBadURL(t, func() error {
		_, err := client.Location.Get(id)
		return err
	})
}

func TestLocationList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("location.list.json", t)
	defer data.Close()

	const (
		page   = 1
		region = "Maryland"
	)
	mux.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		if v := r.FormValue("region"); v != region {
			t.Fatalf("Request.FormValue region = %v, want %v", v, region)
		}
		// TODO: check more request query values
		io.Copy(w, data)
	})

	ll, err := client.Location.List(&LocationListRequest{Page: page, Region: region})
	if err != nil {
		t.Fatal(err)
	}
	if len(ll.Locations) <= 0 {
		t.Fatal("Expected >0 locations")
	}

	for _, l := range ll.Locations {
		if n := 6; n != len(l.ID) {
			t.Fatalf("Location ID len = %d, wanted %d", len(l.ID), n)
		}
		if l.Latitude == 0.0 {
			t.Fatal("Expected non-zero latitude")
		}
		if l.Longitude == 0.0 {
			t.Fatal("Expected non-zero longitude")
		}
	}

	testBadURL(t, func() error {
		_, err := client.Location.List(&LocationListRequest{Page: page, Region: region})
		return err
	})
}

func makeTestLocation() *Location {
	return &Location{
		ID:                       "z9H6Hj",
		Name:                     "Bethesda",
		StreetAddress:            "7900 Norfolk Ave",
		Locality:                 "Bethesda",
		Region:                   "Maryland",
		PostalCode:               "20814",
		Phone:                    "301-652-1311",
		Website:                  "http://www.rockbottom.com/bethesda",
		HoursOfOperationExplicit: []string{"Sunday - Thursday: 11am - 1am", "Friday - Saturday: 11am - 2am"},
		Latitude:                 38.988988,
		Longitude:                -77.097413,
		IsPrimary:                "N",
		InPlanning:               "N",
		IsClosed:                 "N",
		OpenToPublic:             "Y",
		LocationType:             "brewpub",
		LocationTypeDisplay:      "Brewpub",
		CountryISOCode:           "US",
		Country: Country{
			ISOCode:     "US",
			Name:        "UNITED STATES",
			DisplayName: "United States",
			ISOThree:    "USA",
			NumberCode:  840,
		},
		YearOpened: "1980",
		BreweryID:  "D1UQzj",
		Brewery:    Brewery{},
	}

}

func TestLocationUpdate(t *testing.T) {
	setup()
	defer teardown()

	location := makeTestLocation()

	mux.HandleFunc("/location/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		checkURLSuffix(t, r, location.ID)

		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}

		checkPostFormValue(t, r, "name", location.Name)
		checkPostFormValue(t, r, "streetAddress", location.StreetAddress)
		checkPostFormValue(t, r, "locality", location.Locality)
		checkPostFormValue(t, r, "region", location.Region)
		checkPostFormValue(t, r, "postalCode", location.PostalCode)
		checkPostFormValue(t, r, "phone", location.Phone)
		checkPostFormValue(t, r, "website", location.Website)
		checkPostFormValue(t, r, "hoursOfOperationExplicit", location.HoursOfOperationExplicit[0])
		checkPostFormValue(t, r, "latitude", fmt.Sprintf("%f", location.Latitude))
		checkPostFormValue(t, r, "longitude", fmt.Sprintf("%f", location.Longitude))
		checkPostFormValue(t, r, "isPrimary", location.IsPrimary)
		checkPostFormValue(t, r, "inPlanning", location.InPlanning)
		checkPostFormValue(t, r, "isClosed", location.IsClosed)
		checkPostFormValue(t, r, "openToPublic", location.OpenToPublic)
		checkPostFormValue(t, r, "locationType", string(location.LocationType))
		checkPostFormValue(t, r, "countryIsoCode", location.CountryISOCode)

		// Check that fields tagged with "-" or "omitempty" are NOT encoded
		checkPostFormDNE(t, r, "id", "ID", "extendedAddress",
			"ExtendedAddress", "hoursOfOperation", "hoursOfOperationNotes", "tourInfo",
			"LocationTypeDisplay", "country", "Country", "yearClosed",
			"breweryID", "BreweryID", "brewery", "Brewery",
			"status", "Status")
	})

	if err := client.Location.Update(location.ID, location); err != nil {
		t.Fatal(err)
	}

	if client.Location.Update(location.ID, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Location.Update(location.ID, location)
	})
}

func TestLocationDelete(t *testing.T) {
	setup()
	defer teardown()

	const id = "z9H6HJ"
	mux.HandleFunc("/location/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "location" {
			t.Fatal("bad URL, expected \"/location/:locationId\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Location ID", http.StatusNotFound)
		}

	})

	if err := client.Location.Delete(id); err != nil {
		t.Fatal(err)
	}

	if err := client.Location.Delete("******"); err == nil {
		t.Fatal("expected HTTP 404")
	}

	testBadURL(t, func() error {
		return client.Location.Delete(id)
	})
}
