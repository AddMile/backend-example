package user

import (
	"context"

	"github.com/AddMile/backend/internal/config"
	"github.com/AddMile/backend/internal/shared/event"

	pubsubkit "github.com/AddMile/backend/internal/kit/pubsub"
)

type Queue struct {
	cfg *config.Config
	pub *pubsubkit.Publisher
}

func NewQueue(cfg *config.Config, pub *pubsubkit.Publisher) *Queue {
	return &Queue{
		cfg: cfg,
		pub: pub,
	}
}

func (q *Queue) PublishUserCreated(ctx context.Context, event event.UserCreatedEvent) error {
	return q.pub.Publish(ctx, q.cfg.TopicUserCreated, event)
}
