package jikan

import (
	"context"
	"net/url"
)

type SeasonsEndpoints service

// GetNow will return the list of anime airing in the current season
// To filter results, pass in query parameters. Refer to documentation for accepted parameters.
//
// https://docs.api.jikan.moe/#/seasons/getseasonnow
func (s *SeasonsEndpoints) GetNow(ctx context.Context, query *url.Values) (*PaginatedResponseBody[Anime], *Response, error) {
	info := new(PaginatedResponseBody[Anime])

	path := "/v4/seasons/now"
	path += "?" + query.Encode()

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
