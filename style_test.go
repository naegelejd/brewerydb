package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestStyleList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/style.list.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	const page = 1
	mux.HandleFunc("/styles/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		io.Copy(w, data)
	})

	sl, err := client.Style.List(page)
	if err != nil {
		t.Fatal(err)
	}
	if len(sl.Styles) <= 0 {
		t.Fatal("Expected >0 styles")
	}

	for _, s := range sl.Styles {
		if s.ID <= 0 {
			t.Fatal("Expected non-zero style ID")
		}
		if s.Name == "" {
			t.Fatal("Expected non-empty style name")
		}
	}
}
