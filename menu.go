package brewerydb

// MenuService provides access to the BreweryDB Menu API.
// Use Client.Menu.
type MenuService struct {
	c *Client
}

// Styles provides a listing of all Beer Styles.
func (ms *MenuService) Styles() ([]Style, error) {
	req, err := ms.c.NewRequest("GET", "/menu/styles", nil)
	if err != nil {
		return nil, err
	}

	res := struct {
		Status  string
		Data    []Style
		Message string
	}{}
	err = ms.c.Do(req, &res)
	return res.Data, err
}

// Categories provides a listing of all Beer Categories.
func (ms *MenuService) Categories() ([]Category, error) {
	req, err := ms.c.NewRequest("GET", "/menu/categories", nil)
	if err != nil {
		return nil, err
	}

	res := struct {
		Status  string
		Data    []Category
		Message string
	}{}
	err = ms.c.Do(req, &res)
	return res.Data, err
}

// Glassware provides a listing of all Beer Glasses.
func (ms *MenuService) Glassware() ([]Glass, error) {
	req, err := ms.c.NewRequest("GET", "/menu/glassware", nil)
	if err != nil {
		return nil, err
	}

	res := struct {
		Status  string
		Data    []Glass
		Message string
	}{}
	err = ms.c.Do(req, &res)
	return res.Data, err
}

// SRM provides a listing of all SRMs.
func (ms *MenuService) SRM() ([]SRM, error) {
	req, err := ms.c.NewRequest("GET", "/menu/srm", nil)
	if err != nil {
		return nil, err
	}

	res := struct {
		Status  string
		Data    []SRM
		Message string
	}{}
	err = ms.c.Do(req, &res)
	return res.Data, err
}

// BeerAvailability provides a listing of all possible Availability states.
func (ms *MenuService) BeerAvailability() ([]Availability, error) {
	req, err := ms.c.NewRequest("GET", "/menu/beer-availability", nil)
	if err != nil {
		return nil, err
	}

	resp := struct {
		Status  string
		Data    []Availability
		Message string
	}{}
	err = ms.c.Do(req, &resp)
	return resp.Data, err
}

// Fluidsize provides a listing of all fluidsizes.
func (ms *MenuService) Fluidsize() ([]Fluidsize, error) {
	req, err := ms.c.NewRequest("GET", "/menu/fluidsize", nil)
	if err != nil {
		return nil, err
	}

	res := struct {
		Status  string
		Data    []Fluidsize
		Message string
	}{}
	err = ms.c.Do(req, &res)
	return res.Data, err
}

// BeerTemperature provides a mapping of BeerTemperatures to their respective descriptions.
func (ms *MenuService) BeerTemperature() (map[BeerTemperature]string, error) {
	req, err := ms.c.NewRequest("GET", "/menu/beer-temperature", nil)
	if err != nil {
		return nil, err
	}

	res := struct {
		Status  string
		Data    map[BeerTemperature]string
		Message string
	}{}
	err = ms.c.Do(req, &res)
	return res.Data, err
}

// Countries provides a listing of all Countries on Earth.
func (ms *MenuService) Countries() ([]Country, error) {
	req, err := ms.c.NewRequest("GET", "/menu/countries", nil)
	if err != nil {
		return nil, err
	}

	res := struct {
		Status  string
		Data    []Country
		Message string
	}{}
	err = ms.c.Do(req, &res)
	return res.Data, err
}

// Ingredients provides a listing of all Ingredients.
func (ms *MenuService) Ingredients() ([]Ingredient, error) {
	req, err := ms.c.NewRequest("GET", "/menu/ingredients", nil)
	if err != nil {
		return nil, err
	}

	res := struct {
		Status  string
		Data    []Ingredient
		Message string
	}{}
	err = ms.c.Do(req, &res)
	return res.Data, err
}

// LocationTypes provides a mapping of LocationTypes to their respective descriptions.
func (ms *MenuService) LocationTypes() (map[LocationType]string, error) {
	req, err := ms.c.NewRequest("GET", "/menu/location-types", nil)
	if err != nil {
		return nil, err
	}

	res := struct {
		Status  string
		Data    map[LocationType]string
		Message string
	}{}
	err = ms.c.Do(req, &res)
	return res.Data, err
}

// FluidsizeVolume provides a mapping of Volumes to their respective descriptions.
func (ms *MenuService) FluidsizeVolume() (map[Volume]string, error) {
	req, err := ms.c.NewRequest("GET", "/menu/fluidsize-volume", nil)
	if err != nil {
		return nil, err
	}

	res := struct {
		Status  string
		Data    map[Volume]string
		Message string
	}{}
	err = ms.c.Do(req, &res)
	return res.Data, err
}

// EventTypes provides a mapping of EventTypes to their respective descriptions.
func (ms *MenuService) EventTypes() (map[EventType]string, error) {
	req, err := ms.c.NewRequest("GET", "/menu/event-types", nil)
	if err != nil {
		return nil, err
	}

	res := struct {
		Status  string
		Data    map[EventType]string
		Message string
	}{}
	err = ms.c.Do(req, &res)
	return res.Data, err
}
