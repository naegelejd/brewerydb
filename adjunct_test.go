package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestAdjunctList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/adjunct.list.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	const page = 1
	mux.HandleFunc("/adjuncts/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		io.Copy(w, data)
	})

	al, err := client.Adjunct.List(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(al.Adjuncts) <= 0 {
		t.Fatal("Expected >0 adjuncts")
	}

	c := "misc"
	for _, a := range al.Adjuncts {
		if c != a.Category {
			t.Fatalf("Adjunct Category = %s, wanted %s", a.Category, c)
		}
	}
}
