package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/AddMile/backend/internal"
	"github.com/AddMile/backend/internal/shared"
	"github.com/AddMile/backend/internal/shared/event"
)

type UpsertUserParams struct {
	Email    string          `validate:"required,email"`
	Language shared.Language `validate:"required,oneof=en es"`
}

func (s *Service) UpsertUser(ctx context.Context, params UpsertUserParams) (uuid.UUID, error) {
	s.logger.Info("upsert user", slog.Any("params", params))

	if err := s.validator.Validate(params); err != nil {
		s.logger.Error("upsert user: validation failed", slog.Any("params", params), slog.Any("error", err))

		return uuid.Nil, fmt.Errorf("%w: %w", internal.ErrValidation, err)
	}

	data := UpsertUserData(params)
	userID, err := s.storage.UpsertUser(ctx, data)
	if err != nil {
		s.logger.Error("cannot upsert user", slog.Any("data", data), slog.Any("error", err))

		return uuid.Nil, fmt.Errorf("cannot upsert user: %w", err)
	}

	err = s.queue.PublishUserCreated(ctx, event.UserCreatedEvent{UserID: userID})
	if err != nil {
		s.logger.Error("cannot publish user created event", slog.Any("userId", userID), slog.Any("error", err))

		return uuid.Nil, fmt.Errorf("cannot publish user created event: %w", err)
	}

	return userID, nil
}

type UserParams struct {
	UserID uuid.UUID `validate:"required,uuid"`
}

func (s *Service) User(ctx context.Context, params UserParams) (User, error) {
	s.logger.Debug("user by id", slog.Any("params", params))

	data := UserData(params)
	u, err := s.storage.User(ctx, data)
	if err != nil {
		s.logger.Error("cannot get user", slog.Any("data", data), slog.Any("error", err))

		return User{}, fmt.Errorf("cannot get user: %w", err)
	}

	return u, nil
}
