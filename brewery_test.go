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

func TestBreweryGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("brewery.get.json", t)
	defer data.Close()

	const id = "jmGoBA"
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, id)
		io.Copy(w, data)
	})

	b, err := client.Brewery.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if b.ID != id {
		t.Fatalf("Brewery ID = %v, want %v", b.ID, id)
	}

	testBadURL(t, func() error {
		_, err := client.Brewery.Get(id)
		return err
	})
}

func TestBreweryList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("brewery.list.json", t)
	defer data.Close()

	const established = "1988"
	mux.HandleFunc("/breweries", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		if v := r.FormValue("established"); v != established {
			t.Fatalf("Request.FormValue established = %v, wanted %v", v, established)
		}
		// TODO: check more request query values
		io.Copy(w, data)
	})

	bl, err := client.Brewery.List(&BreweryListRequest{Established: established})
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Breweries) <= 0 {
		t.Fatal("Expected >0 breweries")
	}

	for _, b := range bl.Breweries {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Brewery ID len = %d, wanted %d", len(b.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Brewery.List(&BreweryListRequest{Established: established})
		return err
	})
}

func makeTestBrewery() *Brewery {
	return &Brewery{
		ID:             "jmGoBA",
		Name:           "Flying Dog Brewery",
		Description:    "Good people drink good beer.",
		MailingListURL: "boss@flyingdogales.com",
		Image:          "https://s3.amazonaws.com/brewerydbapi/brewery/jmGoBA/upload_0z9L4W-large.png",
		Established:    "1983",
		IsOrganic:      "N",
		Website:        "http://www.flyingdogales.com",
		Status:         "verified",
	}
}

func TestBreweryAdd(t *testing.T) {
	setup()
	defer teardown()

	brewery := makeTestBrewery()

	const newID = "abcdef"
	mux.HandleFunc("/breweries", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")

		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}

		checkPostFormValue(t, r, "name", brewery.Name)
		checkPostFormValue(t, r, "description", brewery.Description)
		checkPostFormValue(t, r, "mailingListUrl", brewery.MailingListURL)
		checkPostFormValue(t, r, "image", brewery.Image)
		checkPostFormValue(t, r, "established", brewery.Established)
		checkPostFormValue(t, r, "isOrganic", brewery.IsOrganic)
		checkPostFormValue(t, r, "website", brewery.Website)

		// Check that fields tagged with "-" or "omitempty" are NOT encoded
		checkPostFormDNE(t, r, "id", "ID", "status", "Status")

		fmt.Fprintf(w, `{"status":"...", "data":{"id":"%s"}, "message":"..."}`, newID)
	})

	id, err := client.Brewery.Add(brewery)
	if err != nil {
		t.Fatal(err)
	}
	if id != newID {
		t.Fatalf("new Brewery ID = %v, want %v", id, newID)
	}

	_, err = client.Brewery.Add(nil)
	if err == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		_, err = client.Brewery.Add(brewery)
		return err
	})
}

func TestBreweryUpdate(t *testing.T) {
	setup()
	defer teardown()

	brewery := makeTestBrewery()

	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		checkURLSuffix(t, r, brewery.ID)

		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}

		checkPostFormValue(t, r, "name", brewery.Name)
		checkPostFormValue(t, r, "description", brewery.Description)
		checkPostFormValue(t, r, "mailingListUrl", brewery.MailingListURL)
		checkPostFormValue(t, r, "image", brewery.Image)
		checkPostFormValue(t, r, "established", brewery.Established)
		checkPostFormValue(t, r, "isOrganic", brewery.IsOrganic)
		checkPostFormValue(t, r, "website", brewery.Website)

		// Check that fields tagged with "-" or "omitempty" are NOT encoded
		checkPostFormDNE(t, r, "id", "ID", "status", "Status")
	})

	if err := client.Brewery.Update(brewery.ID, brewery); err != nil {
		t.Fatal(err)
	}

	if client.Brewery.Update(brewery.ID, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Brewery.Update(brewery.ID, brewery)
	})
}

func TestBreweryDelete(t *testing.T) {
	setup()
	defer teardown()

	const id = "jmGoBA"
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "brewery" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}
	})

	if err := client.Brewery.Delete(id); err != nil {
		t.Fatal(err)
	}

	if err := client.Brewery.Delete("******"); err == nil {
		t.Fatal("expected HTTP 404")
	}

	testBadURL(t, func() error {
		return client.Brewery.Delete(id)
	})
}

func TestBreweryGetSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("brewery.get.socialaccount.json", t)
	defer data.Close()

	const (
		breweryID       = "jmGoBA"
		socialAccountID = 16
	)
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccount" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/socialaccount/:socialaccountId\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(socialAccountID) {
			http.Error(w, "invalid SocialAccount ID", http.StatusNotFound)
		}
		io.Copy(w, data)

	})

	a, err := client.Brewery.GetSocialAccount(breweryID, socialAccountID)
	if err != nil {
		t.Fatal(err)
	}

	if a.ID != socialAccountID {
		t.Fatalf("SocialAccount ID = %v, want %v", a.ID, socialAccountID)
	}

	testBadURL(t, func() error {
		_, err := client.Brewery.GetSocialAccount(breweryID, socialAccountID)
		return err
	})
}

func TestBreweryListSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("brewery.list.socialaccounts.json", t)
	defer data.Close()

	const breweryID = "jmGoBA"
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccounts" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/socialaccounts\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	al, err := client.Brewery.ListSocialAccounts(breweryID)
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
		_, err := client.Brewery.ListSocialAccounts(breweryID)
		return err
	})
}

func TestBreweryAddSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	account := &SocialAccount{
		ID:            3,
		SocialMediaID: 1,
		SocialSite: SocialSite{
			ID:      1,
			Name:    "Facebook Fan Page",
			Website: "http://www.facebook.com",
		},
		Handle: "flying_dog",
	}

	const id = "jmGoBA"
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccounts" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/socialaccounts\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "socialmediaId", strconv.Itoa(account.SocialMediaID))
		checkPostFormValue(t, r, "handle", account.Handle)

		checkPostFormDNE(t, r, "id", "ID", "socialMedia", "SocialSite")
	})

	if err := client.Brewery.AddSocialAccount(id, account); err != nil {
		t.Fatal(err)
	}

	if client.Brewery.AddSocialAccount("******", account) == nil {
		t.Fatal("expected HTTP error")
	}

	if client.Brewery.AddSocialAccount(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Brewery.AddSocialAccount(id, account)
	})
}

func TestBreweryUpdateSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	account := &SocialAccount{
		ID:            3,
		SocialMediaID: 1,
		SocialSite: SocialSite{
			ID:      1,
			Name:    "Facebook Fan Page",
			Website: "http://www.facebook.com",
		},
		Handle: "flying_dog",
	}

	const id = "jmGoBA"
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccount" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/socialaccount/:socialaccountId\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(account.ID) {
			http.Error(w, "invalid SocialAccount ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "socialmediaId", strconv.Itoa(account.SocialMediaID))
		checkPostFormValue(t, r, "handle", account.Handle)

		checkPostFormDNE(t, r, "id", "socialMedia", "SocialSite")
	})

	if err := client.Brewery.UpdateSocialAccount(id, account); err != nil {
		t.Fatal(err)
	}

	if client.Brewery.UpdateSocialAccount("******", account) == nil {
		t.Fatal("expected HTTP error")
	}

	if client.Brewery.UpdateSocialAccount(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Brewery.UpdateSocialAccount(id, account)
	})
}

func TestAddAlternateName(t *testing.T) {
	setup()
	defer teardown()

	const (
		breweryID = "jmGoBA"
		altName   = "Flying Dog"
		newID     = 3
	)
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "brewery" || split[3] != "alternatenames" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/alternatenames\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "name", altName)

		fmt.Fprintf(w, `{"status":"...", "data":{"id":%d}, "message":"..."}`, newID)
	})

	id, err := client.Brewery.AddAlternateName(breweryID, altName)
	if err != nil {
		t.Fatal(err)
	}
	if id != newID {
		t.Fatalf("alternate name ID = %v, want %v", id, newID)
	}

	_, err = client.Brewery.AddAlternateName("******", altName)
	if err == nil {
		t.Fatal("expected HTTP 404")
	}

	testBadURL(t, func() error {
		_, err := client.Brewery.AddAlternateName(breweryID, altName)
		return err
	})
}

func TestBreweryDeleteAlternatName(t *testing.T) {
	setup()
	defer teardown()

	const (
		breweryID = "jmGoBA"
		altID     = 2
	)
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "brewery" || split[3] != "alternatename" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/alternatename/:alternatenameId\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(altID) {
			http.Error(w, "invalid alternatename ID", http.StatusNotFound)
		}
	})

	if err := client.Brewery.DeleteAlternateName(breweryID, altID); err != nil {
		t.Fatal(err)
	}

	if err := client.Brewery.DeleteAlternateName("******", altID); err == nil {
		t.Fatal("expected HTTP 404")
	}

	if err := client.Brewery.DeleteAlternateName(breweryID, -1); err == nil {
		t.Fatal("expected HTTP 404")
	}

	testBadURL(t, func() error {
		return client.Brewery.DeleteAlternateName(breweryID, altID)
	})
}

func TestBreweryAddGuild(t *testing.T) {
	setup()
	defer teardown()

	const (
		breweryID = "jmGoBA"
		guildID   = "k2jMtH"
	)
	discount := "10%"
	firstTest := true
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "brewery" || split[3] != "guilds" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/guilds\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "guildId", guildID)
		if firstTest {
			checkPostFormValue(t, r, "discount", discount)
		}
	})

	if err := client.Brewery.AddGuild(breweryID, guildID, &discount); err != nil {
		t.Fatal(err)
	}

	firstTest = false
	if err := client.Brewery.AddGuild("******", guildID, nil); err == nil {
		t.Fatal("expected HTTP 404")
	}

	if err := client.Brewery.AddGuild(breweryID, guildID, nil); err != nil {
		t.Fatal(err)
	}

	testBadURL(t, func() error {
		return client.Brewery.AddGuild(breweryID, guildID, &discount)
	})
}

func TestBreweryDeleteGuild(t *testing.T) {
	setup()
	defer teardown()

	const (
		breweryID = "jmGoBA"
		guildID   = "k2jMtH"
	)
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "brewery" || split[3] != "guild" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/guild/:guildId\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}
		if split[4] != guildID {
			http.Error(w, "invalid Guild ID", http.StatusNotFound)
		}
	})

	if err := client.Brewery.DeleteGuild(breweryID, guildID); err != nil {
		t.Fatal(err)
	}

	if err := client.Brewery.DeleteGuild("******", guildID); err == nil {
		t.Fatal("expected HTTP 404")
	}

	if err := client.Brewery.DeleteGuild(breweryID, "~~~~~~"); err == nil {
		t.Fatal("expected HTTP 404")
	}

	testBadURL(t, func() error {
		return client.Brewery.DeleteGuild(breweryID, guildID)
	})
}

func TestBreweryAddLocation(t *testing.T) {
	setup()
	defer teardown()

	location := makeTestLocation()

	const (
		breweryID = "jmGoBA"
		newID     = "abcdef"
	)
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")

		split := strings.Split(r.URL.Path, "/")
		if split[1] != "brewery" || split[3] != "locations" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/locations\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}

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

		fmt.Fprintf(w, `{"status":"...", "data":{"guid":"%s"}, "message":"..."}`, newID)
	})

	id, err := client.Brewery.AddLocation(breweryID, location)
	if err != nil {
		t.Fatal(err)
	}
	if id != newID {
		t.Fatalf("Location ID = %v, want %v", id, newID)
	}

	_, err = client.Brewery.AddLocation("******", location)
	if err == nil {
		t.Fatal("expected HTTP 404 error")
	}

	_, err = client.Brewery.AddLocation(breweryID, nil)
	if err == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		_, err := client.Brewery.AddLocation(breweryID, location)
		return err
	})
}

func TestBreweryDeleteSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	const (
		breweryID = "jmGoBA"
		socialID  = 2
	)
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "brewery" || split[3] != "socialaccount" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/socialaccount/:socialaccountId\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(socialID) {
			http.Error(w, "invalid socialaccount ID", http.StatusNotFound)
		}
	})

	if err := client.Brewery.DeleteSocialAccount(breweryID, socialID); err != nil {
		t.Fatal(err)
	}

	if err := client.Brewery.DeleteSocialAccount("******", socialID); err == nil {
		t.Fatal("expected HTTP 404")
	}

	if err := client.Brewery.DeleteSocialAccount(breweryID, -1); err == nil {
		t.Fatal("expected HTTP 404")
	}

	testBadURL(t, func() error {
		return client.Brewery.DeleteSocialAccount(breweryID, socialID)
	})
}

func TestBreweryListEvents(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("event.list.json", t)
	defer data.Close()

	const breweryID = "jmGoBA"
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "events" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/events\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid brewery ID", http.StatusNotFound)
		}

		checkFormValue(t, r, "onlyWinners", "Y")

		io.Copy(w, data)

	})

	el, err := client.Brewery.ListEvents(breweryID, true)
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
		_, err := client.Brewery.ListEvents(breweryID, false)
		return err
	})
}

func TestBreweryListBeers(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("beer.list.json", t)
	defer data.Close()

	const breweryID = "o9TSOv"
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "beers" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/beers\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}

		checkFormValue(t, r, "withBreweries", "Y")
		checkFormValue(t, r, "withSocialAccounts", "Y")
		checkFormValue(t, r, "withIngredients", "Y")

		io.Copy(w, data)
	})

	req := &BreweryBeersRequest{
		WithBreweries:      "Y",
		WithSocialAccounts: "Y",
		WithIngredients:    "Y",
	}
	bl, err := client.Brewery.ListBeers(breweryID, req)
	if err != nil {
		t.Fatal(err)
	}

	if len(bl) <= 0 {
		t.Fatal("Expected >0 Beers")
	}

	for _, b := range bl {
		if l := 6; l != len(b.ID) {
			t.Fatalf("Brewery ID len = %d, wanted %d", len(b.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Brewery.ListBeers(breweryID, req)
		return err
	})
}

func TestBreweryListGuilds(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("guild.list.json", t)
	defer data.Close()

	const breweryID = "o9TSOv"
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "guilds" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/guilds\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	ll, err := client.Brewery.ListGuilds(breweryID)
	if err != nil {
		t.Fatal(err)
		t.Fatal(err)
	}

	if len(ll) <= 0 {
		t.Fatal("Expected >0 Guilds")
	}

	for _, loc := range ll {
		if l := 6; l != len(loc.ID) {
			t.Fatalf("Brewery ID len = %d, wanted %d", len(loc.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Brewery.ListGuilds(breweryID)
		return err
	})
}

func TestBreweryListLocations(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("location.list.json", t)
	defer data.Close()

	const breweryID = "o9TSOv"
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "locations" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/locations\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid Brewery ID", http.StatusNotFound)
		}

		io.Copy(w, data)
	})

	ll, err := client.Brewery.ListLocations(breweryID)
	if err != nil {
		t.Fatal(err)
	}

	if len(ll) <= 0 {
		t.Fatal("Expected >0 Locations")
	}

	for _, loc := range ll {
		if l := 6; l != len(loc.ID) {
			t.Fatalf("Brewery ID len = %d, wanted %d", len(loc.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Brewery.ListLocations(breweryID)
		return err
	})
}

func TestBreweryListAlternateNames(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("brewery.list.alternatenames.json", t)
	defer data.Close()

	const breweryID = "tNDKBY"
	mux.HandleFunc("/brewery/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "alternatenames" {
			t.Fatal("bad URL, expected \"/brewery/:breweryId/alternatenames\"")
		}
		if split[2] != breweryID {
			http.Error(w, "invalid brewery ID", http.StatusNotFound)
		}

		io.Copy(w, data)

	})

	al, err := client.Brewery.ListAlternateNames(breweryID)
	if err != nil {
		t.Fatal(err)
	}

	if len(al) <= 0 {
		t.Fatal("Expected >0 AlternateNames")
	}

	for _, alt := range al {
		if alt.ID <= 0 {
			t.Fatalf("Expected ID >0")
		}
	}

	testBadURL(t, func() error {
		_, err := client.Brewery.ListAlternateNames(breweryID)
		return err
	})
}

func TestBreweryGetRandom(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/brewery.get.random.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/brewery/random", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")

		// TODO: check more request query values
		checkFormValue(t, r, "established", "1983")

		io.Copy(w, data)
	})

	b, err := client.Brewery.GetRandom(&RandomBreweryRequest{Established: "1983"})
	if err != nil {
		t.Fatal(err)
	}

	// Can't really verify specific information since it's a random brewery
	if len(b.Name) <= 0 {
		t.Fatal("Expected non-empty brewery name")
	}
	if len(b.ID) <= 0 {
		t.Fatal("Expected non-empty brewery ID")
	}

	testBadURL(t, func() error {
		_, err := client.Brewery.GetRandom(&RandomBreweryRequest{Established: "1983"})
		return err
	})
}

func ExampleBreweryService_List() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	// Get all breweries established in 1983
	bl, err := c.Brewery.List(&BreweryListRequest{Established: "1983"})
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range bl.Breweries {
		fmt.Println(b.Name, b.ID)
	}

	// Get all information about brewery with given ID (Flying Dog)
	b, err := c.Brewery.Get("jmGoBA")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b.Name)
	fmt.Println(b.Description)
}
