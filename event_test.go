package brewerydb

import (
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestEventGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("event.get.json", t)
	defer data.Close()

	const id = "0oZVAo"
	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, id)
		io.Copy(w, data)
	})

	e, err := client.Event.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if e.ID != id {
		t.Fatalf("Event ID = %v, want %v", e.ID, id)
	}
}

func TestEventList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("event.list.json", t)
	defer data.Close()

	const year = 2015
	mux.HandleFunc("/events/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		if v := r.FormValue("year"); v != strconv.Itoa(year) {
			t.Fatalf("Request.FormValue year = %v, wanted %v", v, year)
		}
		// TODO: check more request query values
		io.Copy(w, data)
	})

	el, err := client.Event.List(&EventListRequest{Year: year})
	if err != nil {
		t.Fatal(err)
	}
	if len(el.Events) <= 0 {
		t.Fatal("Expected >0 events")
	}
	for _, e := range el.Events {
		if l := 6; l != len(e.ID) {
			t.Fatalf("Event ID len = %d, wanted %d", len(e.ID), l)
		}
	}
}
