package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
	v := url.Values{}
	v.Set("p", strconv.Itoa(page))
	u := is.c.url("/ingredients", &v)
	var resp *http.Response
	resp, err = is.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get ingredients")
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&il); err != nil {
		return
	}

	return
}

// Get returns the Ingredient with the given Ingredient ID.
func (is *IngredientService) Get(id int) (ing Ingredient, err error) {
	// GET: /ingredient/:ingredientId
	u := is.c.url(fmt.Sprintf("/ingredient/%d", id), nil)
	var resp *http.Response
	resp, err = is.c.Get(u)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to get ingredient")
		return
	}
	defer resp.Body.Close()

	ingredientResponse := struct {
		Status  string
		Data    Ingredient
		Message string
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&ingredientResponse); err != nil {
		return
	}
	ing = ingredientResponse.Data

	return
}
