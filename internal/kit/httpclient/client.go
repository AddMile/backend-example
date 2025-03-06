package httpclient

import (
	"net/http"
	"time"
)

type DefaultHeaders map[string]string

type HTTPClient struct {
	*http.Client
	defaultHeaders DefaultHeaders
}

// defaultTimeout defines global request timeout.
const defaultTimeout = 20 * time.Second

// HTTPClientOption represents an option for configuring the HTTP client.
type HTTPClientOption func(*HTTPClient) error

func WithTimeout(timeout time.Duration) HTTPClientOption {
	return func(c *HTTPClient) error {
		c.Timeout = timeout

		return nil
	}
}

func WithDefaultHeaders(headers DefaultHeaders) HTTPClientOption {
	return func(c *HTTPClient) error {
		c.defaultHeaders = headers

		return nil
	}
}

// newHTTPClient creates a new HTTP client with the provided options.
func New(opts ...HTTPClientOption) (*HTTPClient, error) {
	httpClient := &http.Client{
		Timeout: defaultTimeout,
	}

	c := &HTTPClient{
		Client: httpClient,
	}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}
