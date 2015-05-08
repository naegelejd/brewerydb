package brewerydb

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestDeleteBeer(t *testing.T) {
	// Attempt to delete non-existent beer
	// err := c.Beer.Delete("zzzzzzzzzzzzzzzzzz")
	// if err == nil {
	// 	t.Fatal("successfully delete a non-existent beer")
	// }
	// t.Fatal(err)
}

// Get a random beer with an ABV between 8.0 and 9.0
func ExampleBeerService_Random() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	req := &RandomBeerRequest{
		ABV: "8",
	}
	b, err := c.Beer.Random(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(b.Name)
	fmt.Println(b.Style.Name)
	fmt.Println(b.Labels.Large)
}

func ExampleBeerService_List() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	// Get first 40 beers with an ABV between 8.0 and 9.0, descending, alphabetical
	bl, err := c.Beer.List(&BeerListRequest{ABV: "8", Sort: SortDescending})
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range bl.Beers {
		fmt.Println(b.Name, b.ID)
	}
}

func ExampleBeerService_Breweries() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

	// Get breweries for a given beer
	breweries, err := c.Beer.Breweries("jmGoBA")
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range breweries {
		fmt.Println(b.Name)
	}
}
