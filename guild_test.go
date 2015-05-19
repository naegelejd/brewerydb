package brewerydb

import (
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
	mux.HandleFunc("/guilds/", func(w http.ResponseWriter, r *http.Request) {
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
}
