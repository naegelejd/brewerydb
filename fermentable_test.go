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
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		io.Copy(w, data)
	})

	fl, err := client.Fermentable.List(1)
	if err != nil {
		t.Error(err)
	}
	if len(fl.Fermentables) <= 0 {
		t.Error("Expected >0 fermentables")
	}
}