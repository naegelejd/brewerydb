package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestIngredientList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/ingredient.list.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
	defer data.Close()

	const page = 1
	mux.HandleFunc("/ingredients/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		io.Copy(w, data)
	})

	il, err := client.Ingredient.List(page)
	if err != nil {
		t.Fatal(err)
	}
	if len(il.Ingredients) <= 0 {
		t.Fatal("Expected >0 ingredients")
	}

	for _, i := range il.Ingredients {
		if i.ID <= 0 {
			t.Fatal("Expected non-zero ingredient ID")
		}
		if i.Category == "" {
			t.Fatal("Expected non-empty ingredient Category")
		}
	}
}
