package jikan

import (
	"net/http"
	"testing"
)

func TestGetAnimeFullById(t *testing.T) {
	client := NewJikanClient()

	_, resp, err := client.Anime.GetAnimeFullById(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.Status)
	}
}

func TestGetAnimeById(t *testing.T) {
	client := NewJikanClient()

	_, resp, err := client.Anime.GetAnimeById(t.Context(), "1")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.Status)
	}
}
