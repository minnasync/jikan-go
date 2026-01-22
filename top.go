package jikan

import (
	"context"
	"net/url"
	"strconv"
	"time"
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

	if s.client.cache != nil {
		err := s.client.cache.Get(ctx, "jikan:top-anime"+query.Encode(), info)
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

	resp, err := s.client.Do(ctx, req, info)
	if err != nil {
		return nil, &Response{
			IsCached: false,
			Response: resp,
		}, err
	}

	if s.client.cache != nil {
		s.client.cache.DeferSet(ctx, "jikan:top-anime"+query.Encode(), info, time.Hour*24)

		animeMap := make(map[string]any, len(info.Data))
		for _, anime := range info.Data {
			animeMap["jikan:anime:"+strconv.Itoa(anime.MalID)] = anime
		}
		s.client.cache.DeferBulkSet(ctx, animeMap, time.Hour*24)
	}

	return info, &Response{
		IsCached: false,
		Response: resp,
	}, nil
}
