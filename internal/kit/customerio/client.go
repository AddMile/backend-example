package customerio

import (
	"context"
	"fmt"
	"time"

	analytics "github.com/customerio/cdp-analytics-go"
)

type Config struct {
	APIKey    string
	Endpoint  string
	BatchSize int
	Interval  time.Duration
	Verbose   bool
}

type Client struct {
	client analytics.Client
}

func NewClient(config Config) (*Client, error) {
	c, err := analytics.NewWithConfig(config.APIKey, analytics.Config{
		Endpoint:  config.Endpoint,
		Interval:  config.Interval,
		BatchSize: config.BatchSize,
		Verbose:   config.Verbose,
	})
	if err != nil {
		return nil, err
	}

	return &Client{client: c}, nil
}

func (c *Client) Close(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("cannot close customer.io client: %w", ctx.Err())
	default:
	}

	return c.client.Close()
}

func (c *Client) Identity(ctx context.Context, user User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("cannot validate user: %w", err)
	}

	if err := c.client.Enqueue(user.toUser()); err != nil {
		return fmt.Errorf("cannot enqueue user: %w", err)
	}

	return nil
}

func (c *Client) TrackEvent(ctx context.Context, event Event) error {
	if err := event.Validate(); err != nil {
		return fmt.Errorf("cannot validate event: %w", err)
	}

	if err := c.client.Enqueue(event.toAnalyticsEvent()); err != nil {
		return fmt.Errorf("cannot enqueue event: %w", err)
	}

	return nil
}
