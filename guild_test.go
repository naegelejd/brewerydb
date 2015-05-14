package brewerydb

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestGuildList(t *testing.T) {
	setup()
	defer teardown()

	data, err := os.Open("test_data/guild.list.json")
	if err != nil {
		t.Fatal("Failed to open test data file")
	}
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
