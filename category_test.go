package brewerydb

import (
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestCategoryGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("category.get.json", t)
	defer data.Close()

	const id = 3
	mux.HandleFunc("/category/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, strconv.Itoa(id))
		io.Copy(w, data)
	})

	c, err := client.Category.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if c.ID != id {
		t.Fatalf("Category ID = %v, want %v", c.ID, id)
	}
}

func TestCategorylist(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("category.list.json", t)
	defer data.Close()

	mux.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	cl, err := client.Category.List()
	if err != nil {
		t.Error(err)
	}
	if len(cl) <= 0 {
		t.Error("Expected >0 categories")
	}
	for _, c := range cl {
		if c.ID <= 0 {
			t.Fatal("Expected non-zero category ID")
		}
		if c.Name == "" {
			t.Fatal("Expected non-empty category name")
		}
	}
}
