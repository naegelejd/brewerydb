package brewerydb

import (
	"io"
	"net/http"
	"os"
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

	mux.HandleFunc("/features/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	fl, err := client.Feature.List(&FeatureRequest{Page: 1, Year: 2015})
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
