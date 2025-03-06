package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *HTTPClient) GET(ctx context.Context, uri string, response any, opts ...RequestOption) error {
	req, err := c.buildRequest(ctx, http.MethodGet, uri, nil, opts...)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.call(req, response)
}

func (c *HTTPClient) POST(ctx context.Context, uri string, request, response any,
	opts ...RequestOption) error {
	req, err := c.buildRequest(ctx, http.MethodPost, uri, request, opts...)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.call(req, response)
}

func (c *HTTPClient) PUT(ctx context.Context, uri string, request, response any,
	opts ...RequestOption) error {
	req, err := c.buildRequest(ctx, http.MethodPut, uri, request, opts...)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.call(req, response)
}

func (c *HTTPClient) DELETE(ctx context.Context, uri string, response any, opts ...RequestOption) error {
	req, err := c.buildRequest(ctx, http.MethodDelete, uri, nil, opts...)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.call(req, response)
}

// buildRequest builds an HTTP request with the provided parameters.
func (c *HTTPClient) buildRequest(ctx context.Context, method, url string, request any,
	opts ...RequestOption) (*http.Request, error) {
	var requestBody []byte
	if request != nil {
		var err error
		requestBody, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
	}

	rp, err := newRequestParams(opts...)
	if err != nil {
		return nil, fmt.Errorf("building request parameters: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("creating request with context: %w", err)
	}

	if len(c.defaultHeaders) > 0 {
		for key, value := range c.defaultHeaders {
			req.Header.Set(key, value)
		}
	}

	for key, value := range rp.headers {
		req.Header.Set(key, value)
	}

	if requestBody == nil {
		req.Header.Set("Accept", rp.contentType)
	} else {
		req.Header.Set("Content-Type", rp.contentType)
	}

	return req, nil
}

// call executes an HTTP request and handles the response.
func (c *HTTPClient) call(request *http.Request, response any) error {
	res, err := c.Do(request)
	if err != nil {
		return fmt.Errorf("doing request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return APIError{
			Status:  res.StatusCode,
			Message: string(body),
		}
	}

	if response != nil {
		if err := json.Unmarshal(body, response); err != nil {
			return fmt.Errorf("unmarshaling response: %w", err)
		}
	}

	return nil
}
