package jikan

import (
	"net/http"
	"testing"
)

func TestGetAnimeFullById(t *testing.T) {
	client := NewJikanClient()

	_, resp, err := client.Anime.GetFullById(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}

	if resp.Response != nil && resp.Response.StatusCode != http.StatusOK {
		t.Fatal(resp.Response.Status)
	}
}

func TestGetAnimeById(t *testing.T) {
	client := NewJikanClient()

	_, resp, err := client.Anime.GetById(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}

	if resp.Response != nil && resp.Response.StatusCode != http.StatusOK {
		t.Fatal(resp.Response.Status)
	}
}
