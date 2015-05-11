package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestFluidsizeList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/fluidsize.list.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/fluidsizes/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	fl, err := client.Fluidsize.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(fl) <= 0 {
		t.Fatal("Expected >0 fluidsizes")
	}

	for _, f := range fl {
		if f.ID <= 0 {
			t.Fatalf("Expected non-zero ID")
		}
	}
}
