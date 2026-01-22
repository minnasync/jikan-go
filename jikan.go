package jikan

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/minnasync/jikan-go/internal/httpx"
	"github.com/redis/go-redis/v9"
)

type Response struct {
	IsCached bool
	Response *http.Response
}

type Client struct {
	client  *http.Client
	baseUrl *url.URL

	cache Cache

	common  service
	Anime   *AnimeEndpoints
	Seasons *SeasonsEndpoints
	Top     *TopEndpoints
}

type service struct {
	client *Client
}

type ClientOption func(*Client)

// WithRedisCache will enable redis caching.
func WithRedisCache(client *redis.Client) ClientOption {
	return func(c *Client) {
		c.cache = NewRedisCache(client)
	}
}

// WithCache will enable caching with a custom cache manager.
func WithCache(cache Cache) ClientOption {
	return func(c *Client) {
		c.cache = cache
	}
}

func (c *Client) newClient() *Client {
	c.common.client = c

	c.Anime = (*AnimeEndpoints)(&c.common)
	c.Seasons = (*SeasonsEndpoints)(&c.common)
	c.Top = (*TopEndpoints)(&c.common)

	return c
}

// NewJikanClient will create a new, default client.
func NewJikanClient(options ...ClientOption) *Client {
	ratelimit := &httpx.RequestLimitRoundTripper{
		RoundTripper: http.DefaultTransport,
	}

	c := &Client{
		client: &http.Client{
			Transport: ratelimit,
		},
		baseUrl: &url.URL{
			Scheme: "https",
			Host:   "api.jikan.moe",
		},
	}

	for _, option := range options {
		option(c)
	}

	return c.newClient()
}

// NewGETRequest will create a new GET request only.
//
// Jikan only supports GET requests, refer to documentation.
// https://docs.api.jikan.moe/#/section/information/allowed-http(s)-requests
func (c *Client) NewGETRequest(path string) (*http.Request, error) {
	u, err := c.baseUrl.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Do will execute an HTTP request.
func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return nil, err
		}
	}

	return resp, nil
}
