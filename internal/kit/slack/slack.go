package slack

import (
	"context"
	"fmt"

	"github.com/AddMile/backend/internal/kit/httpclient"
)

type Client struct {
	client *httpclient.HTTPClient
}

func NewClient() (*Client, error) {
	httpClient, err := httpclient.New()
	if err != nil {
		return nil, fmt.Errorf("initializing httpclient: %w", err)
	}

	return &Client{client: httpClient}, nil
}

func (c Client) SendEvent(ctx context.Context, event Event) error {
	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}

	payload := map[string]string{
		"text": event.Message,
	}

	return c.client.POST(ctx, event.Endpoint, payload, nil)
}
