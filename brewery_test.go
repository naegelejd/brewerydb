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
}

func TestBreweryList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("brewery.list.json", t)
	defer data.Close()

	const established = "1988"
	mux.HandleFunc("/breweries/", func(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println(b.Website)
}
