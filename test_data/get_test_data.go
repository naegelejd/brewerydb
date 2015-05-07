package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/naegelejd/brewerydb"
)

func main() {
	c := brewerydb.NewClient(os.Getenv("BREWERYDB_API_KEY"))

	var test_data = map[string]TestDataGetter{
		"adjunct.list.json":     adjunctList,
		"beer.random.json":      beerRandom,
		"event.list.json":       eventList,
		"fermentable.list.json": fermentableList,
		"search.beer.json":      searchBeer,
		"search.brewery.json":   searchBrewery,
	}

	for filename, action := range test_data {
		if err := getTestData(c, filename, action); err != nil {
			log.Printf("error getting %s: %s\n", filename, err)
		}
	}
}

type TestDataGetter func(*brewerydb.Client) error

func getTestData(c *brewerydb.Client, filename string, action func(c *brewerydb.Client) error) error {
	log.Printf("Creating %s\n", filename)
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

func beerRandom(c *brewerydb.Client) error {
	_, err := c.Beer.Random(&brewerydb.RandomBeerRequest{ABV: "8"})
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

func adjunctList(c *brewerydb.Client) error {
	_, err := c.Adjunct.List(1)
	return err
}

func eventList(c *brewerydb.Client) error {
	_, err := c.Event.List(&brewerydb.EventRequest{Page: 1, Year: 2015})
	return err
}

func fermentableList(c *brewerydb.Client) error {
	_, err := c.Fermentable.List(1)
	return err
}
