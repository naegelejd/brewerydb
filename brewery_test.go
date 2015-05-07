package brewerydb

import (
	"fmt"
	"log"
)

func ExampleBreweryList() {
	c := NewClient("<your API key>")

	// Get all breweries established in 1983
	bl, err := c.Brewery.List(&BreweryListRequest{Established: "1983"})
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range bl.Breweries {
		fmt.Println(b.Name, b.ID)
	}

	// Get all information about brewery with given ID (Flying Dog)
	b, err := c.Brewery.Get("jmGoBA")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b.Name)
	fmt.Println(b.Description)
	fmt.Println(b.Website)
}
