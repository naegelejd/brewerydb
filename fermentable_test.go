package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestFermentableList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/fermentable.list.json")
	if err != nil {
		t.Errorf("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/fermentables/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	fl, err := client.Fermentable.List(1)
	if err != nil {
		t.Error(err)
	}
	if len(fl.Fermentables) <= 0 {
		t.Error("Expected >0 fermentables")
	}
	for _, f := range fl.Fermentables {
		if f.ID <= 0 {
			t.Fatal("Expected non-zero fermentable ID")
		}
		if f.Name == "" {
			t.Fatal("Expected non-empty fermentable name")
		}
	}
}
