package brewerydb

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestFeatureList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/feature.list.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	const (
		page = 1
		year = 2015
	)
	mux.HandleFunc("/features/", func(w http.ResponseWriter, r *http.Request) {
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
}
