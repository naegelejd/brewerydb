package brewerydb

import (
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestSocialSiteGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("socialsite.get.json", t)
	defer data.Close()

	const id = 4
	mux.HandleFunc("/socialsite/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, strconv.Itoa(id))
		io.Copy(w, data)
	})

	s, err := client.SocialSite.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != id {
		t.Fatalf("Socialsite ID = %v, want %v", s.ID, id)
	}
}

func TestSocialSiteList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("socialsite.list.json", t)
	defer data.Close()

	mux.HandleFunc("/socialsites", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	ssl, err := client.SocialSite.List()
	if err != nil {
		t.Error(err)
	}
	if len(ssl) <= 0 {
		t.Error("Expected >0 socialsites")
	}
	for _, s := range ssl {
		if s.ID <= 0 {
			t.Fatal("Expected non-zero socialsite ID")
		}
		if s.Website == "" {
			t.Fatal("Expected non-empty socialsite Website")
		}
	}
}
