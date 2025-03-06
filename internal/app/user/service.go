package user

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/AddMile/backend/internal"
	"github.com/AddMile/backend/internal/shared/event"
)

type Validator interface {
	Validate(any) error
}

type UserStorage interface {
	UpsertUser(context.Context, UpsertUserData) (uuid.UUID, error)
	User(context.Context, UserData) (User, error)
}

type UserQueue interface {
	PublishUserCreated(context.Context, event.UserCreatedEvent) error
}

type Service struct {
	logger    *slog.Logger
	validator Validator
	storage   UserStorage
	queue     UserQueue
}

func NewService(
	logger *slog.Logger,
	validator Validator,
	storage UserStorage,
	queue UserQueue,
) (*Service, error) {
	if logger == nil {
		return nil, internal.ErrNilLogger
	}

	if validator == nil {
		return nil, internal.ErrNilValidator
	}

	if storage == nil {
		return nil, internal.ErrNilStorage
	}

	if queue == nil {
		return nil, internal.ErrNilQueue
	}

	return &Service{
		logger:    logger,
		validator: validator,
		storage:   storage,
		queue:     queue,
	}, nil
}
