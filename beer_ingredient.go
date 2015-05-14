package brewerydb

import "net/http"

// ListIngredients returns a slice of Ingredients found in the Beer with the given ID.
func (bs *BeerService) ListIngredients(beerID string) (el []Ingredient, err error) {
	// GET: /beer/:beerId/ingredients
	var req *http.Request
	req, err = bs.c.NewRequest("GET", "/beer/"+beerID+"/ingredients", nil)
	if err != nil {
		return
	}

	resp := struct {
		Status  string
		Data    []Ingredient
		Message string
	}{}
	err = bs.c.Do(req, &resp)
	return resp.Data, err
}
