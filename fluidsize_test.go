package brewerydb

import (
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestFluidsizeGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("fluidsize.get.json", t)
	defer data.Close()

	const id = 5
	mux.HandleFunc("/fluidsize/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, strconv.Itoa(id))
		io.Copy(w, data)
	})

	f, err := client.Fluidsize.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if f.ID != id {
		t.Fatalf("Fluidsize ID = %v, want %v", f.ID, id)
	}

	testBadURL(t, func() error {
		_, err := client.Fluidsize.Get(id)
		return err
	})
}

func TestFluidsizeList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("fluidsize.list.json", t)
	defer data.Close()

	mux.HandleFunc("/fluidsizes", func(w http.ResponseWriter, r *http.Request) {
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

	testBadURL(t, func() error {
		_, err := client.Fluidsize.List()
		return err
	})
}
