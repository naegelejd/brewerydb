package brewerydb

import (
	"fmt"
	"net/http"
	"testing"
)

var fakeDataChanges = `{
  "status" : "success",
  "numberOfPages" : 5,
  "data" : [
    {
      "action" : "edit",
      "subAttribute" : {
        "locality" : "Some City",
        "locationTypeDisplay" : "Micro Brewery",
        "status" : "verified",
        "statusDisplay" : "Verified",
        "country" : {
          "displayName" : "United States",
          "isoCode" : "US",
          "numberCode" : 840,
          "createDate" : "2012-01-03 02:41:33",
          "name" : "UNITED STATES",
          "isoThree" : "USA"
        },
        "yearOpened" : "2008",
        "updateDate" : "2013-06-21 12:47:17",
        "region" : "Texas",
        "latitude" : 36.863446,
        "inPlanning" : "N",
        "name" : "Main Brewery",
        "hoursOfOperation" : "Visit website for tours and hours.",
        "id" : "qxAQSJ",
        "openToPublic" : "Y",
        "isClosed" : "N",
        "locationType" : "micro",
        "longitude" : -74.676669,
        "phone" : "512-707-2337",
        "website" : "http://www.512brewing.com/",
        "postalCode" : "33333",
        "brewery" : {
          "isOrganic" : "N",
          "website" : "http://www.512brewing.com/",
          "images" : {
            "medium" : "https://s3.amazonaws.com/brewerydbapi/brewery/cJio9R/upload_knal4J-medium.png",
            "large" : "https://s3.amazonaws.com/brewerydbapi/brewery/cJio9R/upload_knal4J-large.png",
            "icon" : "https://s3.amazonaws.com/brewerydbapi/brewery/cJio9R/upload_knal4J-icon.png"
          },
          "id" : "cJio9R",
          "established" : "2008",
          "status" : "verified",
          "updateDate" : "2012-09-29 12:20:39",
          "description" : "(512) Brewing Company is a microbrewery located in the heart of Austin...",
          "statusDisplay" : "Verified",
          "name" : "(512) Brewing Company",
          "createDate" : "2012-01-03 02:41:43"
        },
        "isPrimary" : "Y",
        "countryIsoCode" : "US",
        "createDate" : "2012-01-03 02:41:43",
        "breweryId" : "cJio9R",
        "extendedAddress" : "1000",
        "streetAddress" : "100 Test Street"
      },
      "attribute" : {
        "isOrganic" : "N",
        "website" : "http://www.512brewing.com/",
        "images" : {
          "medium" : "https://s3.amazonaws.com/brewerydbapi/brewery/cJio9R/upload_knal4J-medium.png",
          "large" : "https://s3.amazonaws.com/brewerydbapi/brewery/cJio9R/upload_knal4J-large.png",
          "icon" : "https://s3.amazonaws.com/brewerydbapi/brewery/cJio9R/upload_knal4J-icon.png"
        },
        "id" : "cJio9R",
        "established" : "2008",
        "status" : "verified",
        "updateDate" : "2012-09-29 12:20:39",
        "description" : "(512) Brewing Company is a microbrewery located in the heart of Austin...",
        "statusDisplay" : "Verified",
        "name" : "(512) Brewing Company",
        "createDate" : "2012-01-03 02:41:43"
      },
      "attributeName" : "brewery",
      "subAction" : "edit",
      "changeDate" : "2013-06-24 12:47:17",
      "subAttributeName" : "location"
    }
  ],
  "currentPage" : 3,
  "totalResults" : 215
}`

func TestChangeList(t *testing.T) {
	setup()
	defer teardown()

	const (
		page        = 3
		attrName    = ChangeBrewery
		attrID      = "cJio9R"
		subAttrName = ChangeLocation
		subAttrID   = "qxAQSJ"
	)
	mux.HandleFunc("/changes", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)

		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse URL query", http.StatusBadRequest)
		}

		checkFormValue(t, r, "attributeName", string(attrName))
		checkFormValue(t, r, "attributeId", attrID)

		fmt.Fprint(w, fakeDataChanges)
	})

	req := &ChangeListRequest{
		Page:          page,
		AttributeName: attrName,
		AttributeID:   attrID,
	}

	cl, err := client.Change.List(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(cl.Changes) <= 0 {
		t.Fatal("Expected >0 Changes")
	}

	for _, c := range cl.Changes {
		if c.AttributeName != attrName {
			t.Fatalf("AttributeName = %v, want %v", c.AttributeName, attrName)
		}
		if !testValidChangeAction(c.Action) {
			t.Fatalf("Action = %v, want %v, %v or %v",
				c.Action, ChangeAdd, ChangeDelete, ChangeEdit)
		}
		if l := 6; l != len(c.Attribute.ID) {
			t.Fatalf("Attribute ID len = %v, want %v", len(c.Attribute.ID), l)
		}
		if c.SubAttributeName != subAttrName {
			t.Fatalf("SubAttributeName = %v, want %v", c.SubAttributeName, subAttrName)
		}
		if !testValidChangeAction(c.SubAction) {
			t.Fatalf("SubAction = %v, want %v, %v or %v",
				c.SubAction, ChangeAdd, ChangeDelete, ChangeEdit)
		}
		if l := 6; l != len(c.SubAttribute.ID) {
			t.Fatalf("SubAttribute ID len = %v, want %v", len(c.SubAttribute.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err = client.Change.List(req)
		return err
	})
}

func testValidChangeAction(action ChangeAction) bool {
	switch action {
	case ChangeAdd:
	case ChangeDelete:
	case ChangeEdit:
	default:
		return false
	}
	return true
}
