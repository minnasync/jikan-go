package httpx

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

var (
	// Jikan allows 3 requests per second or 60 request per minute.
	limiter = rate.NewLimiter(
		rate.Every(time.Minute/60),
		3,
	)
)

type RequestLimitRoundTripper struct {
	http.RoundTripper
}

func (r *RequestLimitRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := limiter.Wait(req.Context()); err != nil {
		return nil, err
	}

	return http.DefaultTransport.RoundTrip(req)
}
