package brewerydb

import (
	"encoding/json"
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

// Categories returns all possible Beer Categories.
func (c *CategoryService) Categories() ([]Category, error) {
	// GET: /categories
	u := c.c.url("/categories", nil)

	resp, err := c.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to get categories")
	}
	defer resp.Body.Close()

	categoriesResponse := struct {
		Status  string
		Data    []Category
		Message string
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&categoriesResponse); err != nil {
		return nil, err
	}

	return categoriesResponse.Data, nil
}

// Category obtains the Category with the given Category ID.
func (c *CategoryService) Category(id int) (cat Category, err error) {
	// GET: /category/:categoryId
	u := c.c.url(fmt.Sprintf("/categories/%d", id), nil)

	var resp *http.Response
	resp, err = c.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get category")
		return
	}
	defer resp.Body.Close()

	categoryResponse := struct {
		Status  string
		Data    Category
		Message string
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&categoryResponse); err != nil {
		return
	}

	return categoryResponse.Data, nil
}
