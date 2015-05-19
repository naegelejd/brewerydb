package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/naegelejd/brewerydb"
)

var overwriteFiles bool

func main() {
	var key string
	flag.StringVar(&key, "apikey", "", "brewerydb API key (default: $BREWERYDB_API_KEY)")
	flag.BoolVar(&overwriteFiles, "overwrite", false, "overwrite existing test data")
	flag.Parse()

	if key == "" {
		key = os.Getenv("BREWERYDB_API_KEY")
	}
	c := brewerydb.NewClient(key)

	var test_data = map[string]TestDataGetter{
		"adjunct.get.json":      adjunctGet,
		"adjunct.list.json":     adjunctList,
		"beer.get.json":         beerGet,
		"beer.list.json":        beerList,
		"beer.getrandom.json":   beerGetRandom,
		"brewery.get.json":      breweryGet,
		"brewery.list.json":     breweryList,
		"category.get.json":     categoryGet,
		"category.list.json":    categoryList,
		"event.get.json":        eventGet,
		"event.list.json":       eventList,
		"feature.get.json":      featureGet,
		"feature.list.json":     featureList,
		"feature.byweek.json":   featureByWeek,
		"fluidsize.get.json":    fluidsizeGet,
		"fluidsize.list.json":   fluidsizeList,
		"fermentable.get.json":  fermentableGet,
		"fermentable.list.json": fermentableList,
		"glass.get.json":        glassGet,
		"glass.list.json":       glassList,
		"guild.get.json":        guildGet,
		"guild.list.json":       guildList,
		"hop.get.json":          hopGet,
		"hop.list.json":         hopList,
		"ingredient.get.json":   ingredientGet,
		"ingredient.list.json":  ingredientList,
		"location.get.json":     locationGet,
		"location.list.json":    locationList,
		"search.beer.json":      searchBeer,
		"search.brewery.json":   searchBrewery,
		"search.event.json":     searchEvent,
		"search.guild.json":     searchGuild,
		"socialsite.get.json":   socialsiteGet,
		"socialsite.list.json":  socialsiteList,
		"style.get.json":        styleGet,
		"style.list.json":       styleList,
		"yeast.get.json":        yeastGet,
		"yeast.list.json":       yeastList,
	}

	for filename, action := range test_data {
		if err := getTestData(c, filename, action); err != nil {
			log.Printf("error getting %s: %s\n", filename, err)
		}
	}
}

type TestDataGetter func(*brewerydb.Client) error

func getTestData(c *brewerydb.Client, filename string, action func(c *brewerydb.Client) error) error {
	if _, err := os.Stat(filename); err == nil {
		if overwriteFiles {
			log.Printf("Overwriting %s\n", filename)
		} else {
			log.Printf("Skipping %s\n", filename)
			return nil
		}
	} else {
		log.Printf("Creating %s\n", filename)
	}

	var in, out bytes.Buffer
	c.JSONWriter = &in

	log.Println("Executing API request")
	if err := action(c); err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() {
		f.Close()
		c.JSONWriter = nil
	}()

	if err := json.Indent(&out, in.Bytes(), "", "\t"); err != nil {
		return err
	}

	log.Printf("Saving test data to %s\n", filename)
	_, err = io.Copy(f, &out)
	return err
}

func adjunctGet(c *brewerydb.Client) error {
	_, err := c.Adjunct.Get(923)
	return err
}

func adjunctList(c *brewerydb.Client) error {
	_, err := c.Adjunct.List(1)
	return err
}

func beerGet(c *brewerydb.Client) error {
	_, err := c.Beer.Get("o9TSOv")
	return err
}

func beerList(c *brewerydb.Client) error {
	_, err := c.Beer.List(&brewerydb.BeerListRequest{Page: 1, ABV: "8"})
	return err
}

func beerGetRandom(c *brewerydb.Client) error {
	_, err := c.Beer.GetRandom(&brewerydb.RandomBeerRequest{ABV: "8"})
	return err
}

func breweryGet(c *brewerydb.Client) error {
	_, err := c.Brewery.Get("jmGoBA")
	return err
}

func breweryList(c *brewerydb.Client) error {
	_, err := c.Brewery.List(&brewerydb.BreweryListRequest{Page: 1, Established: "1988"})
	return err
}

func categoryGet(c *brewerydb.Client) error {
	_, err := c.Category.Get(3)
	return err
}

func categoryList(c *brewerydb.Client) error {
	_, err := c.Category.List()
	return err
}

func eventGet(c *brewerydb.Client) error {
	_, err := c.Event.Get("0oZVAo")
	return err
}

func eventList(c *brewerydb.Client) error {
	_, err := c.Event.List(&brewerydb.EventListRequest{Page: 1, Year: 2015})
	return err
}

func featureGet(c *brewerydb.Client) error {
	_, err := c.Feature.Get()
	return err
}

func featureList(c *brewerydb.Client) error {
	_, err := c.Feature.List(&brewerydb.FeatureListRequest{Page: 1, Year: 2015})
	return err
}

func featureByWeek(c *brewerydb.Client) error {
	_, err := c.Feature.ByWeek(2014, 7)
	return err
}

func fermentableGet(c *brewerydb.Client) error {
	_, err := c.Fermentable.Get(753)
	return err
}

func fermentableList(c *brewerydb.Client) error {
	_, err := c.Fermentable.List(1)
	return err
}

func fluidsizeGet(c *brewerydb.Client) error {
	_, err := c.Fluidsize.Get(5)
	return err
}

func fluidsizeList(c *brewerydb.Client) error {
	_, err := c.Fluidsize.List()
	return err
}

func glassGet(c *brewerydb.Client) error {
	_, err := c.Glass.Get(7)
	return err
}

func glassList(c *brewerydb.Client) error {
	_, err := c.Glass.List()
	return err
}

func guildGet(c *brewerydb.Client) error {
	_, err := c.Guild.Get("k2jMtH")
	return err
}

func guildList(c *brewerydb.Client) error {
	_, err := c.Guild.List(&brewerydb.GuildListRequest{Page: 1, Name: "Brewers Association of Maryland"})
	return err
}

func hopGet(c *brewerydb.Client) error {
	_, err := c.Hop.Get(42)
	return err
}

func hopList(c *brewerydb.Client) error {
	_, err := c.Hop.List(1)
	return err
}

func ingredientGet(c *brewerydb.Client) error {
	_, err := c.Ingredient.Get(42)
	return err
}

func ingredientList(c *brewerydb.Client) error {
	_, err := c.Ingredient.List(1)
	return err
}

func locationGet(c *brewerydb.Client) error {
	_, err := c.Location.Get("z9H6HJ")
	return err
}

func locationList(c *brewerydb.Client) error {
	_, err := c.Location.List(&brewerydb.LocationListRequest{Page: 1, Region: "Maryland"})
	return err
}

func searchBeer(c *brewerydb.Client) error {
	_, err := c.Search.Beer("flying", &brewerydb.SearchRequest{Page: 1})
	return err
}

func searchBrewery(c *brewerydb.Client) error {
	_, err := c.Search.Brewery("dog", &brewerydb.SearchRequest{Page: 1})
	return err
}

func searchEvent(c *brewerydb.Client) error {
	_, err := c.Search.Event("festival", &brewerydb.SearchRequest{Page: 1})
	return err
}

func searchGuild(c *brewerydb.Client) error {
	_, err := c.Search.Guild("maryland", &brewerydb.SearchRequest{Page: 1})
	return err
}

func socialsiteGet(c *brewerydb.Client) error {
	_, err := c.SocialSite.Get(4)
	return err
}

func socialsiteList(c *brewerydb.Client) error {
	_, err := c.SocialSite.List()
	return err
}

func styleGet(c *brewerydb.Client) error {
	_, err := c.Style.Get(42)
	return err
}

func styleList(c *brewerydb.Client) error {
	_, err := c.Style.List(1)
	return err
}

func yeastGet(c *brewerydb.Client) error {
	_, err := c.Yeast.Get(1835)
	return err
}

func yeastList(c *brewerydb.Client) error {
	_, err := c.Yeast.List(1)
	return err
}
