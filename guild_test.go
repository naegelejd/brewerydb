package brewerydb

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestGuildGet(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("guild.get.json", t)
	defer data.Close()

	const id = "k2jMtH"
	mux.HandleFunc("/guild/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkURLSuffix(t, r, id)
		io.Copy(w, data)
	})

	g, err := client.Guild.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if g.ID != id {
		t.Fatalf("Guild ID = %v, want %v", g.ID, id)
	}

	testBadURL(t, func() error {
		_, err := client.Guild.Get(id)
		return err
	})
}

func TestGuildList(t *testing.T) {
	setup()
	defer teardown()

	data := loadTestData("guild.list.json", t)
	defer data.Close()

	const (
		page = 1
		name = "Brewers Association of Maryland"
	)
	mux.HandleFunc("/guilds", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		checkPage(t, r, page)
		if v := r.FormValue("name"); v != name {
			t.Fatalf("Request.FormValue name = %v, wanted %v", v, name)
			// TODO: check more request query values
		}
		io.Copy(w, data)
	})

	gl, err := client.Guild.List(&GuildListRequest{Page: page, Name: name})
	if err != nil {
		t.Fatal(err)
	}
	if len(gl.Guilds) <= 0 {
		t.Fatal("Expected >0 guilds")
	}

	for _, g := range gl.Guilds {
		if l := 6; l != len(g.ID) {
			t.Fatalf("Guild ID len = %d, wanted %d", len(g.ID), l)
		}
	}

	testBadURL(t, func() error {
		_, err := client.Guild.List(&GuildListRequest{Page: page, Name: name})
		return err
	})
}

func makeTestGuild() *Guild {
	return &Guild{
		ID:          "k2jMtH",
		Name:        "Brewers Association of Maryland",
		Description: "Non-profit trade association",
		Website:     "http://www.MarylandBeer.org/",
		Image:       "https://s3.amazonaws.com/brewerydbapi/guild/k2jMtH/upload_TjDXP0-large.png",
		Images: Images{
			"https://s3.amazonaws.com/brewerydbapi/guild/k2jMtH/upload_TjDXP0-icon.png",
			"https://s3.amazonaws.com/brewerydbapi/guild/k2jMtH/upload_TjDXP0-medium.png",
			"https://s3.amazonaws.com/brewerydbapi/guild/k2jMtH/upload_TjDXP0-large.png",
		},
		Established: 1996,
	}
}

func TestGuildAdd(t *testing.T) {
	setup()
	defer teardown()

	guild := makeTestGuild()

	const newID = "abcdef"
	mux.HandleFunc("/guilds", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")

		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}

		checkPostFormValue(t, r, "name", guild.Name)
		checkPostFormValue(t, r, "description", guild.Description)
		checkPostFormValue(t, r, "website", guild.Website)
		checkPostFormValue(t, r, "image", guild.Image)
		checkPostFormValue(t, r, "established", strconv.Itoa(guild.Established))

		checkPostFormDNE(t, r, "id", "ID", "images", "Images")

		fmt.Fprintf(w, `{"status":"...", "data":{"id":"%s"}, "message":"..."}`, newID)
	})

	id, err := client.Guild.Add(guild)
	if err != nil {
		t.Fatal(err)
	}
	if id != newID {
		t.Fatalf("new Guild ID = %v, want %v", id, newID)
	}

	_, err = client.Guild.Add(nil)
	if err == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		_, err := client.Guild.Add(guild)
		return err
	})
}

func TestGuildUpdate(t *testing.T) {
	setup()
	defer teardown()

	guild := makeTestGuild()

	mux.HandleFunc("/guild/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		checkURLSuffix(t, r, guild.ID)

		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
		}

		checkPostFormValue(t, r, "name", guild.Name)
		checkPostFormValue(t, r, "description", guild.Description)
		checkPostFormValue(t, r, "website", guild.Website)
		checkPostFormValue(t, r, "image", guild.Image)
		checkPostFormValue(t, r, "established", strconv.Itoa(guild.Established))

		// Check that fields tagged with "-" or "omitempty" are NOT encoded
		checkPostFormDNE(t, r, "id", "ID", "images", "Images", "status")
	})

	if err := client.Guild.Update(guild.ID, guild); err != nil {
		t.Fatal(err)
	}

	if client.Guild.Update(guild.ID, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Guild.Update(guild.ID, guild)
	})
}

func TestGuildDelete(t *testing.T) {
	setup()
	defer teardown()

	const id = "k2jMtH"
	mux.HandleFunc("/guild/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "guild" {
			t.Fatal("bad URL, expected \"/guild/:guildId\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Guild ID", http.StatusNotFound)
		}
	})

	if err := client.Guild.Delete(id); err != nil {
		t.Fatal(err)
	}

	if err := client.Guild.Delete("******"); err == nil {
		t.Fatal("expected HTTP 404")
	}

	testBadURL(t, func() error {
		return client.Guild.Delete(id)
	})
}

func TestGuildAddSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	account := &SocialAccount{
		ID:            1,
		SocialMediaID: 2,
		SocialSite: SocialSite{
			ID:      2,
			Name:    "Twitter",
			Website: "https://www.twitter.com",
		},
		Handle: "marylandbeer",
	}

	const id = "k2jMtH"
	mux.HandleFunc("/guild/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "POST")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccounts" {
			t.Fatal("bad URL, expected \"/guild/:guildId/socialaccounts\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Guild ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "socialmediaId", strconv.Itoa(account.SocialMediaID))
		checkPostFormValue(t, r, "handle", account.Handle)

		checkPostFormDNE(t, r, "id", "ID", "socialMedia", "SocialSite")
	})

	if err := client.Guild.AddSocialAccount(id, account); err != nil {
		t.Fatal(err)
	}

	if client.Guild.AddSocialAccount("******", account) == nil {
		t.Fatal("expected HTTP error")
	}

	if client.Guild.AddSocialAccount(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Guild.AddSocialAccount(id, account)
	})
}

func TestGuildUpdateSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	account := &SocialAccount{
		ID:            1,
		SocialMediaID: 2,
		SocialSite: SocialSite{
			ID:      2,
			Name:    "Twitter",
			Website: "https://www.twitter.com",
		},
		Handle: "marylandbeer",
	}

	const id = "k2jMtH"
	mux.HandleFunc("/guild/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")
		split := strings.Split(r.URL.Path, "/")
		if split[3] != "socialaccount" {
			t.Fatal("bad URL, expected \"/guild/:guildId/socialaccount/:socialaccountId\"")
		}
		if split[2] != id {
			http.Error(w, "invalid Guild ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(account.ID) {
			http.Error(w, "invalid SocialAccount ID", http.StatusNotFound)
		}

		checkPostFormValue(t, r, "socialmediaId", strconv.Itoa(account.SocialMediaID))
		checkPostFormValue(t, r, "handle", account.Handle)

		checkPostFormDNE(t, r, "id", "ID", "socialMedia", "SocialSite")
	})

	if err := client.Guild.UpdateSocialAccount(id, account); err != nil {
		t.Fatal(err)
	}

	if client.Guild.UpdateSocialAccount("******", account) == nil {
		t.Fatal("expected HTTP error")
	}

	if client.Guild.UpdateSocialAccount(id, nil) == nil {
		t.Fatal("expected error regarding nil parameter")
	}

	testBadURL(t, func() error {
		return client.Guild.UpdateSocialAccount(id, account)
	})
}

func TestGuildDeleteSocialAccount(t *testing.T) {
	setup()
	defer teardown()

	const (
		guildID  = "k2jMtH"
		socialID = 2
	)
	mux.HandleFunc("/guild/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "DELETE")
		split := strings.Split(r.URL.Path, "/")
		if split[1] != "guild" || split[3] != "socialaccount" {
			t.Fatal("bad URL, expected \"/guild/:guildId/socialaccount/:socialaccountId\"")
		}
		if split[2] != guildID {
			http.Error(w, "invalid Guild ID", http.StatusNotFound)
		}
		if split[4] != strconv.Itoa(socialID) {
			http.Error(w, "invalid socialaccount ID", http.StatusNotFound)
		}
	})

	if err := client.Guild.DeleteSocialAccount(guildID, socialID); err != nil {
		t.Fatal(err)
	}

	if err := client.Guild.DeleteSocialAccount("******", socialID); err == nil {
		t.Fatal("expected HTTP 404")
	}

	if err := client.Guild.DeleteSocialAccount(guildID, -1); err == nil {
		t.Fatal("expected HTTP 404")
	}

	testBadURL(t, func() error {
		return client.Guild.DeleteSocialAccount(guildID, socialID)
	})
}
