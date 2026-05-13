package jikan

import (
	"context"
	"net/url"
)

type TopEndpoints service

// GetTopAnime will return the top anime.
// To filter results, pass in query parameters. Refer to documentation for accepted parameters.
//
// https://docs.api.jikan.moe/#/top/gettopanime
func (s *TopEndpoints) GetTopAnime(ctx context.Context, query *url.Values) (*PaginatedResponseBody[Anime], *Response, error) {
	info := new(PaginatedResponseBody[Anime])

	path := "/v4/top/anime"
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
