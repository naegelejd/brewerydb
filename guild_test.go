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

	mux.HandleFunc("/guilds/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		io.Copy(w, data)
	})

	gl, err := client.Guild.List(&GuildRequest{Page: 1, Name: "Brewers Association of Maryland"})
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
