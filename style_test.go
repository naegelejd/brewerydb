package brewerydb

import (
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestStyleGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("style.get.json", t)
	defer data.Close()

	const id = 42
	mux.HandleFunc("/style/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, strconv.Itoa(id))
		io.Copy(w, data)
	})

	s, err := client.Style.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != id {
		t.Fatalf("Style ID = %v, want %v", s.ID, id)
	}
}

func TestStyleList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("style.list.json", t)
	defer data.Close()

	const page = 1
	mux.HandleFunc("/styles", func(w http.ResponseWriter, r *http.Request) {
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
