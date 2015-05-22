package brewerydb

import (
	"fmt"
	"net/http"
	"testing"
)

var fakeDataConvertID = `{
  "status" : "success",
  "data" : [
    {
      "oldId" : 23,
      "newId" : "DJcbV1"
    },
    {
      "oldId" : 299,
      "newId" : "n75rVD"
    },
    {
      "oldId" : 599,
      "newId" : "CXuX9r"
    }
  ],
  "message" : "Request Successful"
}`

func TestConvertIDs(t *testing.T) {
	setup()
	defer teardown()

	var (
		ids     = []int{23, 299, 599}
		sids    = "23,299,599"
		mapping = map[int]string{
			23:  "DJcbV1",
			299: "n75rVD",
			599: "CXuX9r",
		}
	)
	convertBeer := true
	mux.HandleFunc("/convertid", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")

		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}

		checkPostFormValue(t, r, "ids", sids)

		tp := r.PostFormValue("type")
		if tp != string(ConvertBeer) && tp != string(ConvertBrewery) {
			http.Error(w, "invalid ID type", http.StatusNotFound)
			return
		}

		if convertBeer {
			checkPostFormValue(t, r, "type", string(ConvertBeer))
		} else {
			checkPostFormValue(t, r, "type", string(ConvertBrewery))
		}

		fmt.Fprint(w, fakeDataConvertID)
	})

	newIDs, err := client.ConvertID.ConvertIDs(ConvertBeer, ids...)
	if err != nil {
		t.Fatal(err)
	}
	for oldID, newID := range newIDs {
		if mapping[oldID] != newID {
			t.Fatalf("New ID = %v, want %v", newID, mapping[oldID])
		}
	}

	convertBeer = false
	newIDs, err = client.ConvertID.ConvertIDs(ConvertBrewery, ids...)
	if err != nil {
		t.Fatal(err)
	}
	for oldID, newID := range newIDs {
		if mapping[oldID] != newID {
			t.Fatalf("New ID = %v, want %v", newID, mapping[oldID])
		}
	}

	_, err = client.ConvertID.ConvertIDs("event", ids...)
	if err == nil {
		t.Fatal("Expected HTTP 404 error")
	}

	testBadURL(t, func() error {
		_, err = client.ConvertID.ConvertIDs(ConvertBrewery, ids...)
		return err
	})

}
