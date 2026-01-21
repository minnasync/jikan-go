package jikan

import (
	"context"
	"net/http"
	"strings"
)

type AnimeEndpoints service

type Anime struct {
	MalID           int              `json:"mal_id"`
	URL             string           `json:"url"`
	Images          map[string]Image `json:"images"`
	Trailer         Trailer          `json:"trailer"`
	Approved        bool             `json:"approved"`
	Titles          []Title          `json:"titles"`
	Title           string           `json:"title"`
	TitleEN         string           `json:"title_english"`
	TitleJP         string           `json:"title_japanese"`
	TitleSynonyms   []string         `json:"title_synonyms"`
	Type            string           `json:"type"`
	Source          string           `json:"source"`
	Status          string           `json:"status"`
	Airing          bool             `json:"airing"`
	Aired           AiredInfo        `json:"aired"`
	Duration        string           `json:"duration"`
	Rating          string           `json:"rating"`
	Score           float64          `json:"score"`
	ScoredBy        int              `json:"scored_by"`
	Rank            int              `json:"rank"`
	Popularity      int              `json:"popularity"`
	Members         int              `json:"members"`
	Favorites       int              `json:"favorites"`
	Episodes        int              `json:"episodes"`
	Synopsis        string           `json:"synopsis"`
	Background      string           `json:"background"`
	Season          string           `json:"season"`
	Year            int              `json:"year"`
	Broadcast       BroadcastInfo    `json:"broadcast"`
	Producers       []Entity         `json:"producers"`
	Licensors       []Entity         `json:"licensors"`
	Studios         []Entity         `json:"studios"`
	Genres          []Entity         `json:"genres"`
	ExplicityGenres []Entity         `json:"explicit_genres"`
	Themes          []Entity         `json:"themes"`
	Demographics    []Entity         `json:"demographics"`
}

// IsExplicit will check the rating to determine if the anime is considered explicit.
func (a *Anime) IsExplicit() bool {
	isExplicit := false

	if strings.HasPrefix(a.Rating, "R") {
		isExplicit = true
	}

	return isExplicit
}

type AnimeFull struct {
	Anime

	Relations []Relation `json:"relations"`
	Theme     Theme      `json:"theme"`
	External  []Link     `json:"external"`
	Streaming []Link     `json:"streaming"`
}

// GetFullById returns a complete anime resource.
//
// https://docs.api.jikan.moe/#/anime/getanimefullbyid
func (s *AnimeEndpoints) GetFullById(ctx context.Context, id string) (*AnimeFull, *http.Response, error) {
	path := "/v4/anime/" + id + "/full"
	req, err := s.client.NewGETRequest(path)
	if err != nil {
		return nil, nil, err
	}

	info := new(ResponseBody[AnimeFull])
	resp, err := s.client.Do(ctx, req, info)
	if err != nil {
		return nil, resp, err
	}

	return &info.Data, resp, nil
}

// GetById returns an anime resource.
//
// https://docs.api.jikan.moe/#/anime/getanimebyid
func (s *AnimeEndpoints) GetById(ctx context.Context, id string) (*Anime, *http.Response, error) {
	path := "/v4/anime/" + id
	req, err := s.client.NewGETRequest(path)
	if err != nil {
		return nil, nil, err
	}

	info := new(ResponseBody[Anime])
	resp, err := s.client.Do(ctx, req, info)
	if err != nil {
		return nil, resp, err
	}

	return &info.Data, resp, nil
}
