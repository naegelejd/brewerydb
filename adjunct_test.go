package brewerydb

import (
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestAdjunctGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("adjunct.get.json", t)
	defer data.Close()

	const id = 923
	mux.HandleFunc("/adjunct/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, strconv.Itoa(id))
		io.Copy(w, data)
	})

	a, err := client.Adjunct.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if a.ID != id {
		t.Fatalf("Adjunct ID = %v, want %v", a.ID, id)
	}
}

func TestAdjunctList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("adjunct.list.json", t)
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
