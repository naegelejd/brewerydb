package brewerydb

import (
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestIngredientGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("ingredient.get.json", t)
	defer data.Close()

	const id = 42
	mux.HandleFunc("/ingredient/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, strconv.Itoa(id))
		io.Copy(w, data)
	})

	i, err := client.Ingredient.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if i.ID != id {
		t.Fatalf("Ingredient ID = %v, want %v", i.ID, id)
	}

	testBadURL(t, func() error {
		_, err := client.Ingredient.Get(id)
		return err
	})
}

func TestIngredientList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("ingredient.list.json", t)
	defer data.Close()

	const page = 1
	mux.HandleFunc("/ingredients", func(w http.ResponseWriter, r *http.Request) {
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

	testBadURL(t, func() error {
		_, err := client.Ingredient.List(page)
		return err
	})
}
