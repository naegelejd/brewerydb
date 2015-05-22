package brewerydb

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func testMenuHelper(t *testing.T, name string, testMenu func() error) {
	setup()
	defer teardown()

	data := loadTestData("menu."+name+".json", t)
	defer data.Close()

	mux.HandleFunc("/menu/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		split := strings.Split(r.URL.Path, "/")
		if split[2] != name {
			t.Fatalf("bad URL, expected \"/menu/%s\"", name)
		}

		io.Copy(w, data)
	})

	if err := testMenu(); err != nil {
		t.Fatal(err)
	}

	testBadURL(t, func() error {
		return testMenu()
	})
}

func TestMenuStyles(t *testing.T) {
	testMenuHelper(t, "styles", func() error {
		l, err := client.Menu.Styles()
		if err != nil {
			return err
		}

		if len(l) <= 0 {
			return fmt.Errorf("Expected >0 Styles")
		}

		for _, v := range l {
			if v.ID <= 0 {
				return fmt.Errorf("Expected ID >0")
			}
		}
		return nil
	})
}

func TestMenuCategories(t *testing.T) {
	testMenuHelper(t, "categories", func() error {
		l, err := client.Menu.Categories()
		if err != nil {
			return err
		}

		if len(l) <= 0 {
			return fmt.Errorf("Expected >0 Categories")
		}

		for _, v := range l {
			if v.ID <= 0 {
				return fmt.Errorf("Expected ID >0")
			}
		}
		return nil
	})
}

func TestMenuGlassware(t *testing.T) {
	testMenuHelper(t, "glassware", func() error {
		l, err := client.Menu.Glassware()
		if err != nil {
			return err
		}

		if len(l) <= 0 {
			return fmt.Errorf("Expected >0 Glassware")
		}

		for _, v := range l {
			if v.ID <= 0 {
				return fmt.Errorf("Expected ID >0")
			}
		}
		return nil
	})
}

func TestMenuSRM(t *testing.T) {
	testMenuHelper(t, "srm", func() error {
		l, err := client.Menu.SRM()
		if err != nil {
			return err
		}

		if len(l) <= 0 {
			return fmt.Errorf("Expected >0 SRM")
		}

		for _, v := range l {
			if v.ID <= 0 {
				return fmt.Errorf("Expected ID >0")
			}
		}
		return nil
	})
}

func TestMenuBeerAvailability(t *testing.T) {
	testMenuHelper(t, "beer-availability", func() error {
		l, err := client.Menu.BeerAvailability()
		if err != nil {
			return err
		}

		if len(l) <= 0 {
			return fmt.Errorf("Expected >0 Availabilities")
		}

		for _, v := range l {
			if v.ID <= 0 {
				return fmt.Errorf("Expected ID >0")
			}
		}
		return nil
	})
}

func TestMenuFluidsize(t *testing.T) {
	testMenuHelper(t, "fluidsize", func() error {
		l, err := client.Menu.Fluidsize()
		if err != nil {
			return err
		}

		if len(l) <= 0 {
			return fmt.Errorf("Expected >0 Fluidsizes")
		}

		for _, v := range l {
			if v.ID <= 0 {
				return fmt.Errorf("Expected ID >0")
			}
		}
		return nil
	})
}

func TestMenuBeerTemperature(t *testing.T) {
	testMenuHelper(t, "beer-temperature", func() error {
		m, err := client.Menu.BeerTemperature()
		if err != nil {
			return err
		}

		if len(m) <= 0 {
			return fmt.Errorf("Expected >0 BeerTemperatures")
		}

		for k, v := range m {
			if len(k) <= 0 || len(v) <= 0 {
				return fmt.Errorf("Expected non-empty BeerTemperature")
			}
		}
		return nil
	})
}

func TestMenuCountries(t *testing.T) {
	testMenuHelper(t, "countries", func() error {
		l, err := client.Menu.Countries()
		if err != nil {
			return err
		}

		if len(l) <= 0 {
			return fmt.Errorf("Expected >0 Countries")
		}
		return nil
	})
}

func TestMenuIngredients(t *testing.T) {
	testMenuHelper(t, "ingredients", func() error {
		l, err := client.Menu.Ingredients()
		if err != nil {
			return err
		}

		if len(l) <= 0 {
			return fmt.Errorf("Expected >0 Ingredients")
		}

		for _, v := range l {
			if v.ID <= 0 {
				return fmt.Errorf("Expected ID >0")
			}
		}
		return nil
	})
}

func TestMenuLocationTypes(t *testing.T) {
	testMenuHelper(t, "location-types", func() error {
		m, err := client.Menu.LocationTypes()
		if err != nil {
			return err
		}

		if len(m) <= 0 {
			return fmt.Errorf("Expected >0 LocationTypes")
		}

		for k, v := range m {
			if len(k) <= 0 || len(v) <= 0 {
				return fmt.Errorf("Expected non-empty LocationType")
			}
		}
		return nil
	})
}

func TestMenuFluidsizeVolume(t *testing.T) {
	testMenuHelper(t, "fluidsize-volume", func() error {
		m, err := client.Menu.FluidsizeVolume()
		if err != nil {
			return err
		}

		if len(m) <= 0 {
			return fmt.Errorf("Expected >0 FluidsizeVolumes")
		}

		for k, v := range m {
			if len(k) <= 0 || len(v) <= 0 {
				return fmt.Errorf("Expected non-empty FluidsizeVolume")
			}
		}
		return nil
	})
}

func TestMenuEventTypes(t *testing.T) {
	testMenuHelper(t, "event-types", func() error {
		m, err := client.Menu.EventTypes()
		if err != nil {
			return err
		}

		if len(m) <= 0 {
			return fmt.Errorf("Expected >0 EventTypes")
		}

		for k, v := range m {
			if len(k) <= 0 || len(v) <= 0 {
				return fmt.Errorf("Expected non-empty EventType")
			}
		}
		return nil
	})
}
