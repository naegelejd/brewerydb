package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestSocialSiteList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/socialsite.list.json")
	if err != nil {
		t.Errorf("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/socialsites/", func(w http.ResponseWriter, r *http.Request) {
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
