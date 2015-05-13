package brewerydb

import (
	"fmt"
	"net/url"
	"testing"
)

func TestEncodeEmbedded(t *testing.T) {
	t.Skip()
	breweryID := "abcdef"

	type Inner struct {
		AwardCategoryID string `json:"awardCategoryId,omitempty"`
		AwardPlaceID    string `json:"awardPlaceId,omitempty"`
	}
	outer0 := struct {
		BreweryID string `json:"breweryId"`
		Inner
	}{breweryID, Inner{}}

	var v url.Values
	v = encode(outer0)

	if x := v.Get("breweryId"); x != breweryID {
		t.Fatalf("breweryID = %v, wanted %v", x, breweryID)
	}

	categoryID := "ghijkl"
	placeID := "mnopqr"
	outer1 := struct {
		BreweryID string `json:"breweryId"`
		Inner
	}{breweryID, Inner{categoryID, placeID}}
	fmt.Printf("%+v\n", outer1)

	v = encode(outer1)
	if x := v.Get("breweryId"); x != breweryID {
		t.Fatalf("breweryID = %v, wanted %v", x, breweryID)
	}
	if x := v.Get("awardCategoryId"); x != categoryID {
		t.Fatalf("awardCategoryId = %v, wanted %v", x, categoryID)
	}
	if x := v.Get("awardPlaceId"); x != placeID {
		t.Fatalf("awardPlaceId = %v, wanted %v", x, placeID)
	}

}
