package brewerydb

import (
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestGlassGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("glass.get.json", t)
	defer data.Close()

	const id = 7
	mux.HandleFunc("/glass/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, strconv.Itoa(id))
		io.Copy(w, data)
	})

	g, err := client.Glass.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if g.ID != id {
		t.Fatalf("Glass ID = %v, want %v", g.ID, id)
	}
}

func TestGlassList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("glass.list.json", t)
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
