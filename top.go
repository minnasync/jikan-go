package jikan

import (
	"context"
	"net/http"
	"net/url"
)

type TopEndpoints service

// GetTopAnime will return the top anime.
// To filter results, pass in query parameters. Refer to documentation for accepted parameters.
//
// https://docs.api.jikan.moe/#/top/gettopanime
func (s *TopEndpoints) GetTopAnime(ctx context.Context, query *url.Values) (*PaginatedResponseBody[Anime], *http.Response, error) {
	path := "/v4/top/anime"
	path += "?" + query.Encode()

	req, err := s.client.NewGETRequest(path)
	if err != nil {
		return nil, nil, err
	}

	info := new(PaginatedResponseBody[Anime])
	resp, err := s.client.Do(ctx, req, info)
	if err != nil {
		return nil, resp, err
	}

	return info, resp, nil
}
