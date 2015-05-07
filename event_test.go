package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestEventList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/event.list.json")
	if err != nil {
		t.Errorf("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/events/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		io.Copy(w, data)
	})

	el, err := client.Event.List(&EventRequest{Year: 2015})
	if err != nil {
		t.Error(err)
	}
	if len(el.Events) <= 0 {
		t.Error("Expected >0 events")
	}
	for _, e := range el.Events {
		if l := 6; l != len(e.ID) {
			t.Errorf("Event ID len = %d, wanted %d", len(e.ID), l)
		}
	}
}

func TestEventGet(t *testing.T) {
	// TODO: don't skip
	t.Skip()

	setup()
	defer teardown()

	data, err := os.Open("test_data/event.get.json")
	if err != nil {
		t.Errorf("Failed to open test data file")
	}
	defer data.Close()

	mux.HandleFunc("/event/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		io.Copy(w, data)
	})

	e, err := client.Event.Get("mB7srw")
	if err != nil {
		t.Errorf(err.Error())
	}
	if n := "Bare Beach Beer Bash"; n != e.Name {
		t.Error("Event name = %v, wanted %v", e.Name, n)
	}
}
