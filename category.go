package brewerydb

import (
	"fmt"
	"net/http"
)

// CategoryService provides access to the BreweryDB Category API. Use Client.Category.
type CategoryService struct {
	c *Client
}

// Category represents a type of Beer as specified by the Brewers Association
// Style Guidelines.
type Category struct {
	ID          int
	Name        string
	Description string
	CreateDate  string
	UpdateDate  string
}

// List returns all possible Beer Categories.
func (cs *CategoryService) List() ([]Category, error) {
	// GET: /categories
	req, err := cs.c.NewRequest("GET", "/categories", nil)
	if err != nil {
		return nil, err
	}

	resp := struct {
		Status  string
		Data    []Category
		Message string
	}{}
	err = cs.c.Do(req, &resp)
	return resp.Data, err
}

// Get obtains the Category with the given Category ID.
func (cs *CategoryService) Get(id int) (cat Category, err error) {
	// GET: /category/:categoryId
	var req *http.Request
	req, err = cs.c.NewRequest("GET", fmt.Sprintf("/category/%d", id), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    Category
		Message string
	}{}
	err = cs.c.Do(req, &resp)
	return resp.Data, err
}
