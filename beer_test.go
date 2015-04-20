package brewerydb

import (
	"testing"
)

func TestDeleteBeer(t *testing.T) {
	c := NewClient("myfakekey")

	// Attempt to delete non-existent beer
	err := c.DeleteBeer("zzzzzzzzzzzzzzzzzz")
	if err == nil {
		t.Fatal("successfully delete a non-existent beer")
	}
	t.Fatal(err)
}
