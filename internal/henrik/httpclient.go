package henrik

import (
	"fmt"
	"github.com/go-zoox/fetch"
	"net/http"
	"strconv"
	"time"
)

type RateLimitError struct {
	Limit     int
	Remaining int
	Reset     int
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("Rate limit exceeded: Limit=%d, Remaining=%d, Reset=%d", e.Limit, e.Remaining, e.Reset)
}

func getRateLimitHeaders(response *fetch.Response) (*RateLimitError, error) {
	limit := response.Headers.Get("x-ratelimit-limit")
	remaining := response.Headers.Get("x-ratelimit-remaining")
	reset := response.Headers.Get("x-ratelimit-reset")

	if limit == "" || remaining == "" || reset == "" {
		return nil, nil
	}

	// Convert headers to integers
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}

	remainingInt, err := strconv.Atoi(remaining)
	if err != nil {
		return nil, err
	}

	resetInt, err := strconv.Atoi(reset)
	if err != nil {
		return nil, err
	}

	return &RateLimitError{
		Limit:     limitInt,
		Remaining: remainingInt,
		Reset:     resetInt,
	}, nil
}

func fetchWithRetry(url string, headers map[string]string) (*fetch.Response, error) {
	for {
		response, err := fetch.Get(url, &fetch.Config{
			Headers: headers,
		})
		if err != nil {
			return nil, err
		}

		if response.StatusCode() == http.StatusOK {
			return response, nil
		}

		if response.StatusCode() == http.StatusTooManyRequests {
			rateLimitError, rateLimitErr := getRateLimitHeaders(response)
			if rateLimitErr != nil {
				return nil, rateLimitErr
			}

			if rateLimitError != nil && rateLimitError.Remaining == 0 {
				time.Sleep(time.Duration(rateLimitError.Reset) * time.Second)
				continue
			}
		}

		return nil, fmt.Errorf("received status code %d", response.StatusCode())
	}
}
