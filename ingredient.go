package brewerydb

import (
	"fmt"
	"net/http"
)

// IngredientService provides access to the BreweryDB Ingredient API.
// Use Client.Ingredient.
type IngredientService struct {
	c *Client
}

// Ingredient represents a single Beer Ingredient.
type Ingredient struct {
	ID              int
	Name            string
	Category        string
	CategoryDisplay string
	CreateDate      string
	UpdateDate      string
}

// IngredientList represents a single "page" containing a slice of Ingredients.
type IngredientList struct {
	CurrentPage   int
	NumberOfPages int
	TotalResults  int
	Ingredients   []Ingredient `json:"data"`
}

// List returns all Ingredients on the given page.
func (is *IngredientService) List(page int) (il IngredientList, err error) {
	// GET: /ingredients
	var req *http.Request
	req, err = is.c.NewRequest("GET", "/ingredients", &Page{page})
	if err != nil {
		return
	}

	err = is.c.Do(req, &il)
	return
}

// Get returns the Ingredient with the given Ingredient ID.
func (is *IngredientService) Get(id int) (ing Ingredient, err error) {
	// GET: /ingredient/:ingredientId
	var req *http.Request
	req, err = is.c.NewRequest("GET", fmt.Sprintf("/ingredient/%d", id), nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    Ingredient
		Message string
	}{}
	err = is.c.Do(req, &resp)
	return resp.Data, err
}
