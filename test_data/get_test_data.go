package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/naegelejd/brewerydb"
)

var overwriteFiles bool

func main() {
	var key string
	flag.StringVar(&key, "apikey", "", "brewerydb API key (default: $BREWERYDB_API_KEY)")
	flag.BoolVar(&overwriteFiles, "overwrite", false, "overwrite existing test data")
	flag.Parse()

	if key == "" {
		key = os.Getenv("BREWERYDB_API_KEY")
	}
	c := brewerydb.NewClient(key)

	var test_data = map[string]TestDataGetter{
		"adjunct.list.json":     adjunctList,
		"beer.list.json":        beerList,
		"beer.random.json":      beerRandom,
		"brewery.list.json":     breweryList,
		"category.list.json":    categoryList,
		"event.list.json":       eventList,
		"feature.list.json":     featureList,
		"fluidsize.list.json":   fluidsizeList,
		"fermentable.list.json": fermentableList,
		"glass.list.json":       glassList,
		"guild.list.json":       guildList,
		"hop.list.json":         hopList,
		"ingredient.list.json":  ingredientList,
		"location.list.json":    locationList,
		"search.beer.json":      searchBeer,
		"search.brewery.json":   searchBrewery,
		"style.list.json":       styleList,
		"yeast.list.json":       yeastList,
	}

	for filename, action := range test_data {
		if err := getTestData(c, filename, action); err != nil {
			log.Printf("error getting %s: %s\n", filename, err)
		}
	}
}

type TestDataGetter func(*brewerydb.Client) error

func getTestData(c *brewerydb.Client, filename string, action func(c *brewerydb.Client) error) error {
	if _, err := os.Stat(filename); err == nil {
		if overwriteFiles {
			log.Printf("Overwriting %s\n", filename)
		} else {
			log.Printf("Skipping %s\n", filename)
			return nil
		}
	} else {
		log.Printf("Creating %s\n", filename)
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() {
		f.Close()
		c.JSONWriter = nil
	}()

	var in, out bytes.Buffer
	c.JSONWriter = &in

	log.Println("Executing API request")
	if err := action(c); err != nil {
		return err
	}

	if err := json.Indent(&out, in.Bytes(), "", "\t"); err != nil {
		return err
	}

	log.Printf("Saving test data to %s\n", filename)
	_, err = io.Copy(f, &out)
	return err
}

func adjunctList(c *brewerydb.Client) error {
	_, err := c.Adjunct.List(1)
	return err
}

func beerList(c *brewerydb.Client) error {
	_, err := c.Beer.List(&brewerydb.BeerListRequest{Page: 1})
	return err
}

func beerRandom(c *brewerydb.Client) error {
	_, err := c.Beer.Random(&brewerydb.RandomBeerRequest{ABV: "8"})
	return err
}

func breweryList(c *brewerydb.Client) error {
	_, err := c.Brewery.List(&brewerydb.BreweryListRequest{Page: 1})
	return err
}

func categoryList(c *brewerydb.Client) error {
	_, err := c.Category.List()
	return err
}

func eventList(c *brewerydb.Client) error {
	_, err := c.Event.List(&brewerydb.EventRequest{Page: 1, Year: 2015})
	return err
}

func featureList(c *brewerydb.Client) error {
	_, err := c.Feature.List(&brewerydb.FeatureRequest{Page: 1})
	return err
}

func fermentableList(c *brewerydb.Client) error {
	_, err := c.Fermentable.List(1)
	return err
}

func fluidsizeList(c *brewerydb.Client) error {
	_, err := c.Fluidsize.List()
	return err
}

func glassList(c *brewerydb.Client) error {
	_, err := c.Glass.List()
	return err
}

func guildList(c *brewerydb.Client) error {
	_, err := c.Guild.List(&brewerydb.GuildRequest{Page: 1})
	return err
}

func hopList(c *brewerydb.Client) error {
	_, err := c.Hop.List(1)
	return err
}

func ingredientList(c *brewerydb.Client) error {
	_, err := c.Ingredient.List(1)
	return err
}

func locationList(c *brewerydb.Client) error {
	_, err := c.Location.List(&brewerydb.LocationRequest{Page: 1})
	return err
}

func searchBeer(c *brewerydb.Client) error {
	_, err := c.Search.Beer("flying", &brewerydb.SearchRequest{Page: 1})
	return err
}

func searchBrewery(c *brewerydb.Client) error {
	_, err := c.Search.Brewery("dog", &brewerydb.SearchRequest{Page: 1})
	return err
}

func styleList(c *brewerydb.Client) error {
	_, err := c.Style.List(1)
	return err
}

func yeastList(c *brewerydb.Client) error {
	_, err := c.Yeast.List(1)
	return err
}
