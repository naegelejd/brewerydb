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

	categoriesResponse := struct {
		Status  string
		Data    []Category
		Message string
	}{}
	if err := cs.c.Do(req, &categoriesResponse); err != nil {
		return nil, err
	}
	return categoriesResponse.Data, nil
}

// Get obtains the Category with the given Category ID.
func (cs *CategoryService) Get(id int) (cat Category, err error) {
	// GET: /category/:categoryId
	var req *http.Request
	req, err = cs.c.NewRequest("GET", fmt.Sprintf("/category/%d", id), nil)
	if err != nil {
		return
	}

	categoryResponse := struct {
		Status  string
		Data    Category
		Message string
	}{}
	if err = cs.c.Do(req, &categoryResponse); err != nil {
		return
	}
	return categoryResponse.Data, nil
}
