package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"time"

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
		"adjunct.get.json":                 adjunctGet,
		"adjunct.list.json":                adjunctList,
		"beer.get.json":                    beerGet,
		"beer.list.json":                   beerList,
		"beer.list.socialaccounts.json":    beerListSocialAccounts,
		"beer.get.socialaccount.json":      beerGetSocialAccount,
		"beer.get.random.json":             beerGetRandom,
		"brewery.get.json":                 breweryGet,
		"brewery.list.json":                breweryList,
		"brewery.list.alternatenames.json": breweryListAlternateNames,
		"brewery.list.socialaccounts.json": breweryListSocialAccounts,
		"brewery.get.socialaccount.json":   breweryGetSocialAccount,
		"brewery.get.random.json":          breweryGetRandom,
		"category.get.json":                categoryGet,
		"category.list.json":               categoryList,
		"event.get.json":                   eventGet,
		"event.list.json":                  eventList,
		"event.get.awardcategory.json":     eventGetAwardCategory,
		"event.list.awardcategories.json":  eventListAwardCategories,
		"event.get.awardplace.json":        eventGetAwardPlace,
		"event.list.awardplaces.json":      eventListAwardPlaces,
		"feature.get.json":                 featureGet,
		"feature.list.json":                featureList,
		"feature.byweek.json":              featureByWeek,
		"fluidsize.get.json":               fluidsizeGet,
		"fluidsize.list.json":              fluidsizeList,
		"fermentable.get.json":             fermentableGet,
		"fermentable.list.json":            fermentableList,
		"glass.get.json":                   glassGet,
		"glass.list.json":                  glassList,
		"guild.get.json":                   guildGet,
		"guild.list.json":                  guildList,
		"guild.list.socialaccounts.json":   guildListSocialAccounts,
		"guild.get.socialaccount.json":     guildGetSocialAccount,
		"hop.get.json":                     hopGet,
		"hop.list.json":                    hopList,
		"ingredient.get.json":              ingredientGet,
		"ingredient.list.json":             ingredientList,
		"location.get.json":                locationGet,
		"location.list.json":               locationList,
		"menu.beer-availability.json":      menuBeerAvailability,
		"menu.glassware.json":              menuGlassware,
		"menu.fluidsize.json":              menuFluidsize,
		"menu.beer-temperature.json":       menuBeerTemperature,
		"menu.countries.json":              menuCountries,
		"menu.styles.json":                 menuStyles,
		"menu.location-types.json":         menuLocationTypes,
		"menu.fluidsize-volume.json":       menuFluidsizeVolume,
		"menu.event-types.json":            menuEventTypes,
		"menu.ingredients.json":            menuIngredients,
		"menu.categories.json":             menuCategories,
		"menu.srm.json":                    menuSRM,
		"search.beer.json":                 searchBeer,
		"search.brewery.json":              searchBrewery,
		"search.event.json":                searchEvent,
		"search.guild.json":                searchGuild,
		"search.style.json":                searchStyle,
		"search.geopoint.json":             searchGeoPoint,
		"search.upc.json":                  searchUPC,
		"socialsite.get.json":              socialsiteGet,
		"socialsite.list.json":             socialsiteList,
		"style.get.json":                   styleGet,
		"style.list.json":                  styleList,
		"yeast.get.json":                   yeastGet,
		"yeast.list.json":                  yeastList,
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

	log.Println("Sleeping...")
	time.Sleep(200 * time.Millisecond)

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

func beerGetSocialAccount(c *brewerydb.Client) error {
	_, err := c.Beer.GetSocialAccount("MwSypd", 1)
	return err
}

func beerListSocialAccounts(c *brewerydb.Client) error {
	_, err := c.Beer.ListSocialAccounts("MwSypd")
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

func breweryListAlternateNames(c *brewerydb.Client) error {
	_, err := c.Brewery.ListAlternateNames("tNDKBY")
	return err
}

func breweryGetSocialAccount(c *brewerydb.Client) error {
	_, err := c.Brewery.GetSocialAccount("d25euF", 16)
	return err
}

func breweryListSocialAccounts(c *brewerydb.Client) error {
	_, err := c.Brewery.ListSocialAccounts("d25euF")
	return err
}

func breweryGetRandom(c *brewerydb.Client) error {
	_, err := c.Brewery.GetRandom(&brewerydb.RandomBreweryRequest{Established: "1983"})
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

func eventGetAwardCategory(c *brewerydb.Client) error {
	_, err := c.Event.GetAwardCategory("cJio9R", 87)
	return err
}

func eventListAwardCategories(c *brewerydb.Client) error {
	_, err := c.Event.ListAwardCategories("cJio9R")
	return err
}

func eventGetAwardPlace(c *brewerydb.Client) error {
	_, err := c.Event.GetAwardPlace("cJio9R", 3)
	return err
}

func eventListAwardPlaces(c *brewerydb.Client) error {
	_, err := c.Event.ListAwardPlaces("cJio9R")
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

func guildListSocialAccounts(c *brewerydb.Client) error {
	_, err := c.Guild.ListSocialAccounts("cJio9R")
	return err
}

func guildGetSocialAccount(c *brewerydb.Client) error {
	_, err := c.Guild.GetSocialAccount("cJio9R", 4)
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

func menuBeerAvailability(c *brewerydb.Client) error {
	_, err := c.Menu.BeerAvailability()
	return err
}
func menuGlassware(c *brewerydb.Client) error {
	_, err := c.Menu.Glassware()
	return err
}
func menuFluidsize(c *brewerydb.Client) error {
	_, err := c.Menu.Fluidsize()
	return err
}
func menuBeerTemperature(c *brewerydb.Client) error {
	_, err := c.Menu.BeerTemperature()
	return err
}
func menuCountries(c *brewerydb.Client) error {
	_, err := c.Menu.Countries()
	return err
}
func menuStyles(c *brewerydb.Client) error {
	_, err := c.Menu.Styles()
	return err
}
func menuLocationTypes(c *brewerydb.Client) error {
	_, err := c.Menu.LocationTypes()
	return err
}
func menuFluidsizeVolume(c *brewerydb.Client) error {
	_, err := c.Menu.FluidsizeVolume()
	return err
}
func menuEventTypes(c *brewerydb.Client) error {
	_, err := c.Menu.EventTypes()
	return err
}
func menuIngredients(c *brewerydb.Client) error {
	_, err := c.Menu.Ingredients()
	return err
}
func menuCategories(c *brewerydb.Client) error {
	_, err := c.Menu.Categories()
	return err
}
func menuSRM(c *brewerydb.Client) error {
	_, err := c.Menu.SRM()
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

func searchStyle(c *brewerydb.Client) error {
	_, err := c.Search.Style("Pale Ale", true)
	return err
}
func searchGeoPoint(c *brewerydb.Client) error {
	_, err := c.Search.GeoPoint(&brewerydb.GeoPointRequest{Latitude: 35.772096, Longitude: -78.638614})
	return err
}
func searchUPC(c *brewerydb.Client) error {
	_, err := c.Search.UPC(606905008303)
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
