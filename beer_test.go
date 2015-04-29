package brewerydb

import (
	"fmt"
	"log"
	"testing"
)

func TestDeleteBeer(t *testing.T) {
	c := NewClient("myfakekey")

	// Attempt to delete non-existent beer
	err := c.Beer.Delete("zzzzzzzzzzzzzzzzzz")
	if err == nil {
		t.Fatal("successfully delete a non-existent beer")
	}
	t.Fatal(err)
}

// Get a random beer with an ABV between 8.0 and 9.0
func ExampleBeerService_Random() {
	c := NewClient("<your API key>")

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

func ExampleBeerList() {
	c := NewClient("<your API key>")

	// Get first 40 beers with an ABV between 8.0 and 9.0, descending, alphabetical
	beers := c.Beer.NewBeerList(&BeerListRequest{ABV: "8", Sort: DescendingSort})
	count := 0
	for b, err := beers.First(); b != nil; b, err = beers.Next() {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(b.Name, b.ID)
		count++
		if count > 40 {
			break
		}
	}
}

func ExampleBeerService_Breweries() {
	c := NewClient("<your API key>")

	// Get breweries for a given beer
	breweries, err := c.Beer.Breweries("jmGoBA")
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range breweries {
		fmt.Println(b.Name)
	}
}
