package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestGlassList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/glass.list.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/glassware/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	gl, err := client.Glass.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(gl) <= 0 {
		t.Fatal("Expected >0 glasses")
	}

	for _, g := range gl {
		if g.ID <= 0 {
			t.Fatal("Expected non-zero Glass ID")
		}
		if g.Name == "" {
			t.Fatal("Expected non-empty Glass Name")
		}
	}
}
