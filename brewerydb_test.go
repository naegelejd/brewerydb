package brewerydb

import (
	"os"
	"testing"
)

var testAddress = "http://localhost:8080"

func TestMain(m *testing.M) {
	apiURL = testAddress
	os.Exit(m.Run())
}
