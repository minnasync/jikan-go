package jikan

import (
	"context"
	"fmt"
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

// Get will return the list of anime airing for a provided year + season.
//
// https://docs.api.jikan.moe/#/seasons/getseason
func (s *SeasonsEndpoints) Get(ctx context.Context, year int, season string, query *url.Values) (*PaginatedResponseBody[Anime], *Response, error) {
	info := new(PaginatedResponseBody[Anime])

	path := fmt.Sprintf("/v4/seasons/%d/%s", year, season)
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

// GetList will return a list of available seasons.
//
// https://docs.api.jikan.moe/#/seasons/getseasonslist
func (s *SeasonsEndpoints) GetList(ctx context.Context) (*PaginatedResponseBody[Season], *Response, error) {
	info := new(PaginatedResponseBody[Season])
	path := "/v4/seasons"

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

	return info, &Response{
		IsCached: false,
		Response: resp,
	}, nil
}

// GetUpcoming will return the list of anime for the upcoming season.
//
// https://docs.api.jikan.moe/#/seasons/getseasonupcoming
func (s *SeasonsEndpoints) GetUpcoming(ctx context.Context, query *url.Values) (*PaginatedResponseBody[Anime], *Response, error) {
	info := new(PaginatedResponseBody[Anime])

	path := "/v4/seasons/upcoming"
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
