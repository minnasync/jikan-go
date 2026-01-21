package jikan

import (
	"context"
	"net/http"
	"net/url"
)

type SeasonsEndpoints service

// Now will return the list of anime airing in the current season
// To filter results, pass in query parameters. Refer to documentation for accepted parameters.
//
// https://docs.api.jikan.moe/#/seasons/getseasonnow
func (s *SeasonsEndpoints) Now(ctx context.Context, query *url.Values) (*PaginatedResponseBody[Anime], *http.Response, error) {
	path := "/v4/seasons/now"
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
