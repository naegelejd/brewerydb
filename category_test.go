package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestCategorclist(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/category.list.json")
	if err != nil {
		t.Errorf("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	cl, err := client.Category.List()
	if err != nil {
		t.Error(err)
	}
	if len(cl) <= 0 {
		t.Error("Expected >0 categories")
	}
	for _, c := range cl {
		if c.ID <= 0 {
			t.Fatal("Expected non-zero category ID")
		}
		if c.Name == "" {
			t.Fatal("Expected non-empty category name")
		}
	}
}
