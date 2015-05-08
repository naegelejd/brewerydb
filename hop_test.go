package brewerydb

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestHop(t *testing.T) {

}

func ExampleHopService_List() {
	c := NewClient(os.Getenv("BREWERYDB_API_KEY"))

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
