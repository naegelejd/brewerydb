package brewerydb

import (
	"fmt"
	"log"
)

func ExampleBreweryList() {
	c := NewClient("<your API key>")

	// Get all breweries established in 1983
	bs := c.Brewery.NewBreweryList(&BreweryListRequest{Established: "1983"})
	for b, err := bs.First(); b != nil; b, err = bs.Next() {
		if err != nil {
			log.Fatal("error: ", err)
		}
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
