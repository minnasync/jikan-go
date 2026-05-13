package jikan

import (
	"context"
	"net/url"
	"strings"
)

type AnimeEndpoints service

type Images map[string]Image

// GetPoster will get the highest resolution poster image.
func (i Images) GetPoster() *string {
	formats := []string{"webp", "jpg"}

	for _, format := range formats {
		images, ok := i[format]
		if !ok {
			continue
		}

		if images.LargeURL != "" {
			return &images.LargeURL
		}

		if images.ImageURL != "" {
			return &images.ImageURL
		}

		if images.SmallURL != "" {
			return &images.SmallURL
		}
	}

	return nil
}

type Anime struct {
	MalID           int           `json:"mal_id"`
	URL             string        `json:"url"`
	Images          Images        `json:"images"`
	Trailer         Trailer       `json:"trailer"`
	Approved        bool          `json:"approved"`
	Titles          []Title       `json:"titles"`
	Title           string        `json:"title"`
	TitleEN         string        `json:"title_english"`
	TitleJP         string        `json:"title_japanese"`
	TitleSynonyms   []string      `json:"title_synonyms"`
	Type            *string       `json:"type"`
	Source          *string       `json:"source"`
	Episodes        *int          `json:"episodes"`
	Status          *string       `json:"status"`
	Airing          bool          `json:"airing"`
	Aired           AiredInfo     `json:"aired"`
	Duration        *string       `json:"duration"`
	Rating          *string       `json:"rating"`
	Score           *float64      `json:"score"`
	ScoredBy        *int          `json:"scored_by"`
	Rank            *int          `json:"rank"`
	Popularity      *int          `json:"popularity"`
	Members         *int          `json:"members"`
	Favorites       *int          `json:"favorites"`
	Synopsis        *string       `json:"synopsis"`
	Background      *string       `json:"background"`
	Season          *string       `json:"season"`
	Year            *int          `json:"year"`
	Broadcast       BroadcastInfo `json:"broadcast"`
	Producers       []Entity      `json:"producers"`
	Licensors       []Entity      `json:"licensors"`
	Studios         []Entity      `json:"studios"`
	Genres          []Entity      `json:"genres"`
	ExplicityGenres []Entity      `json:"explicit_genres"`
	Themes          []Entity      `json:"themes"`
	Demographics    []Entity      `json:"demographics"`
}

// IsExplicit will check the rating to determine if the anime is considered explicit.
func (a *Anime) IsExplicit() bool {
	isExplicit := false

	if a.Rating != nil {
		if strings.HasPrefix(*a.Rating, "R") {
			isExplicit = true
		}
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
func (s *AnimeEndpoints) GetFullById(ctx context.Context, id string) (*AnimeFull, *Response, error) {
	path := "/v4/anime/" + id + "/full"

	if s.client.cache != nil {
		info, err := s.client.cache.Anime().GetAnimeFull(ctx, id)
		if err == nil {
			return info, &Response{
				IsCached: true,
				Response: nil,
			}, nil
		}
	}

	req, err := s.client.NewGETRequest(path)
	if err != nil {
		return nil, nil, err
	}

	info := new(ResponseBody[AnimeFull])
	resp, err := s.client.Do(ctx, req, info)
	if err != nil {
		return nil, &Response{
			IsCached: false,
			Response: resp,
		}, err
	}

	if s.client.cache != nil {
		go func() {
			_ = s.client.cache.Anime().SetAnimeFull(ctx, info.Data)
		}()
	}

	return &info.Data, &Response{
		IsCached: false,
		Response: resp,
	}, nil
}

// GetById returns an anime resource.
//
// https://docs.api.jikan.moe/#/anime/getanimebyid
func (s *AnimeEndpoints) GetById(ctx context.Context, id string) (*Anime, *Response, error) {
	path := "/v4/anime/" + id

	if s.client.cache != nil {
		info, err := s.client.cache.Anime().GetAnime(ctx, id)
		if err == nil {
			return info, &Response{
				IsCached: true,
				Response: nil,
			}, nil
		}
	}

	req, err := s.client.NewGETRequest(path)
	if err != nil {
		return nil, nil, err
	}

	info := new(ResponseBody[Anime])
	resp, err := s.client.Do(ctx, req, info)
	if err != nil {
		return nil, &Response{
			IsCached: false,
			Response: resp,
		}, err
	}

	if s.client.cache != nil {
		go func() {
			_ = s.client.cache.Anime().SetAnime(ctx, info.Data)
		}()
	}

	return &info.Data, &Response{
		IsCached: false,
		Response: resp,
	}, nil
}

// GetSearch will search for an anime based on a query.
//
// https://docs.api.jikan.moe/#/anime/getanimesearch
func (s *AnimeEndpoints) GetSearch(ctx context.Context, query string, values *url.Values) (*PaginatedResponseBody[Anime], *Response, error) {
	info := new(PaginatedResponseBody[Anime])
	path := "/v4/anime"
	if values == nil {
		values = &url.Values{}
	}

	values.Set("q", query)
	path += "?" + values.Encode()

	req, err := s.client.NewGETRequest(path)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(ctx, req, info)
	if err != nil {
		return nil, &Response{
			IsCached: false,
			Response: resp,
		}, err
	}

	if s.client.cache != nil {
		go func() {
			_ = s.client.cache.Anime().BulkSetAnime(ctx, info.Data)
		}()
	}

	return info, &Response{
		IsCached: false,
		Response: resp,
	}, nil
}
