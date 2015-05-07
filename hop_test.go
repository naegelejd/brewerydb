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
	h, err := c.Hop.Get(84)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", h)

	// Get all types of hops
	hl, err := c.Hop.List(1)
	if err != nil {
		log.Fatal(err)
	}
	for _, h := range hl.Hops {
		fmt.Println(h.Name)
	}
}
