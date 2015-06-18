# brewerydb

brewerydb is a Go library for accessing the [BreweryDB API](http://www.brewerydb.com)

[![GoDoc](https://godoc.org/github.com/naegelejd/brewerydb?status.svg)](https://godoc.org/github.com/naegelejd/brewerydb) [![Build Status](https://travis-ci.org/naegelejd/brewerydb.svg)](https://travis-ci.org/naegelejd/brewerydb)[![Coverage Status](https://coveralls.io/repos/naegelejd/brewerydb/badge.svg?branch=master)](https://coveralls.io/r/naegelejd/brewerydb?branch=master)

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
// Get any random beer
beer, _ := client.Beer.Random(&brewerydb.RandomBeerRequest{ABV: "8"})
fmt.Println(beer.Name, beer.Style.Name)
```

or

```go
// Get all breweries established in 1983
bs, err := client.Brewery.List(&brewerydb.BreweryListRequest{Established: "1983"})
if err != nil {
    panic(err)
}
for _, b := range bs {
    fmt.Println(b.Name, b.Website)
}
```

or

```go
// "What is in Dragon's Milk?"
bl, _ := client.Search.Beer("Dragon's Milk", nil)

var beerID string
for _, beer := range bl.Beers {
    if beer.Name == "Dragon's Milk" {
        beerID = beer.ID
    }
}
if beerID == "" {
    panic("Dragon's Milk not found")
}

ingredients, _ := client.Beer.ListIngredients(beerID)
adjuncts, _ := client.Beer.ListAdjuncts(beerID)
fermentables, _ := client.Beer.ListFermentables(beerID)
hops, _ := client.Beer.ListHops(beerID)
yeasts, _ := client.Beer.ListYeasts(beerID)

fmt.Println("Dragon's Milk:")
fmt.Println("  Ingredients:")
for _, ingredient := range ingredients {
    fmt.Println("    " + ingredient.Name)
}
fmt.Println("\n  Adjuncts:")
for _, adjunct := range adjuncts {
    fmt.Println("    " + adjunct.Name)
}
fmt.Println("  Fermentables:")
for _, fermentable := range fermentables {
    fmt.Println("    " + fermentable.Name)
}
fmt.Println("  Hops:")
for _, hop := range hops {
    fmt.Println("    " + hop.Name)
}
fmt.Println("  Yeasts:")
for _, yeast := range yeasts {
    fmt.Println("    " + yeast.Name)
}
```

## status

This library is under development. Please feel free to suggest design changes or report issues.

## license

This library is distributed under the BSD-style license found in the [LICENSE](https://github.com/naegelejd/brewerydb/blob/master/LICENSE) file.


[![views](https://sourcegraph.com/api/repos/github.com/naegelejd/brewerydb/.counters/views.svg?no-count=1)](https://sourcegraph.com/github.com/naegelejd/brewerydb)[![views 24h](https://sourcegraph.com/api/repos/github.com/naegelejd/brewerydb/.counters/views-24h.svg)](https://sourcegraph.com/github.com/naegelejd/brewerydb)
