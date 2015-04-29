package brewerydb

import (
	"fmt"
	"log"
	"testing"
)

func TestHop(t *testing.T) {

}

func ExampleHopList() {
	c := NewClient("<your API key>")

	// Get a specific variety of hop with a given ID
	h, err := c.Hop.Hop(84)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", h)

	// Get all types of hops
	hs := c.Hop.NewHopList()
	for h, err := hs.First(); h != nil; h, err = hs.Next() {
		if err != nil {
			log.Fatal("error: ", err)
		}
		fmt.Println(h.Name)
	}
}
