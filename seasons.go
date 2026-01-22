package jikan

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

type SeasonsEndpoints service

// Now will return the list of anime airing in the current season
// To filter results, pass in query parameters. Refer to documentation for accepted parameters.
//
// https://docs.api.jikan.moe/#/seasons/getseasonnow
func (s *SeasonsEndpoints) Now(ctx context.Context, query *url.Values) (*PaginatedResponseBody[Anime], *Response, error) {
	info := new(PaginatedResponseBody[Anime])

	path := "/v4/seasons/now"
	path += "?" + query.Encode()

	if s.client.cache != nil {
		err := s.client.cache.Get(ctx, "jikan:seasons-now"+query.Encode(), info)
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
		s.client.cache.DeferSet(ctx, "jikan:seasons-now"+query.Encode(), info, time.Hour*24)
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
