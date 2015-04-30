package brewerydb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MenuService provides access to the BreweryDB Menu API.
// Use Client.Menu.
type MenuService struct {
	c *Client
}

// Styles provides a listing of all Beer Styles.
func (ms *MenuService) Styles() ([]Style, error) {
	u := ms.c.url("/menu/styles", nil)
	resp, err := ms.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	res := struct {
		Status  string
		Data    []Style
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// Categories provides a listing of all Beer Categories.
func (ms *MenuService) Categories() ([]Category, error) {
	u := ms.c.url("/menu/categories", nil)
	resp, err := ms.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	res := struct {
		Status  string
		Data    []Category
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// Glassware provides a listing of all Beer Glasses.
func (ms *MenuService) Glassware() ([]Glass, error) {
	u := ms.c.url("/menu/glassware", nil)
	resp, err := ms.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	res := struct {
		Status  string
		Data    []Glass
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// SRM represents a Standard Reference Method.
type SRM struct {
	Hex  string
	Name string
}

// SRM provides a listing of all SRMs.
func (ms *MenuService) SRM() ([]SRM, error) {
	u := ms.c.url("/menu/srm", nil)
	resp, err := ms.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	res := struct {
		Status  string
		Data    map[string]SRM
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	srm := make([]SRM, 0, len(res.Data))
	for _, v := range res.Data {
		srm = append(srm, v)
	}
	return srm, nil
}

// BeerAvailability provides a listing of all possible Availability states.
func (ms *MenuService) BeerAvailability() ([]Availability, error) {
	u := ms.c.url("/menu/beer-availability", nil)
	resp, err := ms.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	res := struct {
		Status  string
		Data    map[string]Availability
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	avail := make([]Availability, 0, len(res.Data))
	for _, v := range res.Data {
		avail = append(avail, v)
	}
	return avail, nil
}

// Fluidsize provides a listing of all fluidsizes.
func (ms *MenuService) Fluidsize() ([]Fluidsize, error) {
	u := ms.c.url("/menu/fluidsize", nil)
	resp, err := ms.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	res := struct {
		Status  string
		Data    []Fluidsize
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// BeerTemperature provides a mapping of BeerTemperatures to their respective descriptions.
func (ms *MenuService) BeerTemperature() (map[BeerTemperature]string, error) {
	u := ms.c.url("/menu/beer-temperature", nil)
	resp, err := ms.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	res := struct {
		Status  string
		Data    map[BeerTemperature]string
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// Countries provides a listing of all Countries on Earth.
func (ms *MenuService) Countries() ([]Country, error) {
	u := ms.c.url("/menu/countries", nil)
	resp, err := ms.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	res := struct {
		Status  string
		Data    []Country
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// Ingredients provides a listing of all Ingredients.
func (ms *MenuService) Ingredients() ([]Ingredient, error) {
	u := ms.c.url("/menu/ingredients", nil)
	resp, err := ms.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	res := struct {
		Status  string
		Data    []Ingredient
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// LocationTypes provides a mapping of LocationTypes to their respective descriptions.
func (ms *MenuService) LocationTypes() (map[LocationType]string, error) {
	u := ms.c.url("/menu/location-types", nil)
	resp, err := ms.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	res := struct {
		Status  string
		Data    map[LocationType]string
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// FluidsizeVolume provides a mapping of Volumes to their respective descriptions.
func (ms *MenuService) FluidsizeVolume() (map[Volume]string, error) {
	u := ms.c.url("/menu/fluidsize-volume", nil)
	resp, err := ms.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	res := struct {
		Status  string
		Data    map[Volume]string
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// EventTypes provides a mapping of EventTypes to their respective descriptions.
func (ms *MenuService) EventTypes() (map[EventType]string, error) {
	u := ms.c.url("/menu/event-types", nil)
	resp, err := ms.c.Get(u)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	res := struct {
		Status  string
		Data    map[EventType]string
		Message string
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return res.Data, nil
}
