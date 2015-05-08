# brewerydb

brewerydb is a library for accessing the [BreweryDB API](http://www.brewerydb.com)

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

## status

This library is under heavy development. Please feel free to suggest design changes or report issues.

## license

This library is distributed under the BSD-style license found in the [LICENSE](https://github.com/naegelejd/brewerydb/blob/master/LICENSE) file.
