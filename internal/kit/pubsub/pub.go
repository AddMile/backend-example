package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type Publisher struct {
	client *pubsub.Client
}

func NewPublisher(ctx context.Context, projectID string) (*Publisher, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("publisher: failed to initialize publisher: %w", err)
	}

	return &Publisher{
		client: client,
	}, nil
}

func (p *Publisher) Close() error {
	return p.client.Close()
}

func (p *Publisher) Publish(ctx context.Context, topic string, message any) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("publisher: cannot serialize message to JSON: %w", err)
	}

	t := p.client.Topic(topic)

	_, err = t.Publish(ctx, &pubsub.Message{Data: jsonData}).Get(ctx)
	if err != nil {
		return fmt.Errorf("publisher: cannot publish message: %w", err)
	}

	return nil
}
