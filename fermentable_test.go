package brewerydb

import (
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestFermentableGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("fermentable.get.json", t)
	defer data.Close()

	const id = 753
	mux.HandleFunc("/fermentable/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, strconv.Itoa(id))
		io.Copy(w, data)
	})

	f, err := client.Fermentable.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if f.ID != id {
		t.Fatalf("Fermentable ID = %v, want %v", f.ID, id)
	}
}

func TestFermentableList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("fermentable.list.json", t)
	defer data.Close()

	const page = 1
	mux.HandleFunc("/fermentables/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		io.Copy(w, data)
	})

	fl, err := client.Fermentable.List(page)
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
