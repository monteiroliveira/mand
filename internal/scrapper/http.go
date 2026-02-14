package scrapper

import (
	"context"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type HttpClient struct {
	client *http.Client
}

func NewHttpClient() *HttpClient {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	retryClient.CheckRetry = checkRetry
	retryClient.Backoff = backoff

	return &HttpClient{
		client: retryClient.StandardClient(),
	}
}

func backoff(_, _ time.Duration, attemptNum int, resp *http.Response) time.Duration {
	mult := math.Pow(2, float64(attemptNum)) * float64(time.Second) * 0.5
	return time.Duration(mult)
}

func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if resp != nil {
		if resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode == http.StatusBadGateway ||
			resp.StatusCode == http.StatusServiceUnavailable {
			return true, nil
		}
	}
	return false, nil
}

func (h *HttpClient) Get(ctx context.Context, url string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// TODO: process custom headers

	response, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
