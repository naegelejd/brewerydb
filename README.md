# brewerydb

brewerydb is a library for accessing the [BreweryDB API](http://www.brewerydb.com)

Documentation: [![GoDoc](https://godoc.org/github.com/naegelejd/brewerydb?status.svg)](https://godoc.org/github.com/naegelejd/brewerydb)

## usage

```go
import "github.com/naegelejd/brewerydb"
```

Construct a new `Client` using your BreweryDB API key:

```go
client := brewerydb.NewClient("<your API key>")
```

Then use the available services to access the API.
For example:

```go
# Get any random beer
beer, _ := client.Beer.Random(nil)
fmt.Println(beer.Name, beer.Style.Name)
```

or

```go
# Get all breweries established in 1983
bs := c.Brewery.NewBreweryList(&brewerydb.BreweryListRequest{Established: "1983"})
for b, err := bs.First(); b != nil; b, err = bs.Next() {
    if err != nil {
        panic(err)
    }
    fmt.Println(b.Name, b.ID)
}
```

# status

This library is still under heavy development. Please feel free to suggest design changes or submit bug fixes.

## license

This library is distributed under the BSD-style license found in the [LICENSE](https://github.com/naegelejd/brewerydb/blob/master/LICENSE) file.
