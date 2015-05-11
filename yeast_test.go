package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestYeastList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/yeast.list.json")
	if err != nil {
		t.Errorf("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/yeasts/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	yl, err := client.Yeast.List(1)
	if err != nil {
		t.Error(err)
	}
	if len(yl.Yeasts) <= 0 {
		t.Error("Expected >0 yeasts")
	}
	for _, y := range yl.Yeasts {
		if y.ID <= 0 {
			t.Fatal("Expected non-zero yeast ID")
		}
		if y.Name == "" {
			t.Fatal("Expected non-empty yeast name")
		}
	}
}
