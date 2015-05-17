package brewerydb

import (
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestYeastGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("yeast.get.json", t)
	defer data.Close()

	const id = 1835
	mux.HandleFunc("/yeast/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, strconv.Itoa(id))
		io.Copy(w, data)
	})

	y, err := client.Yeast.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if y.ID != id {
		t.Fatalf("Yeast ID = %v, want %v", y.ID, id)
	}
}

func TestYeastList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("yeast.list.json", t)
	defer data.Close()

	const page = 1
	mux.HandleFunc("/yeasts/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		io.Copy(w, data)
	})

	yl, err := client.Yeast.List(page)
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
