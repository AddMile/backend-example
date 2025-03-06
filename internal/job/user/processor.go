package user

import (
	"context"
	"log/slog"

	"github.com/AddMile/backend/internal/app/user"
	"github.com/AddMile/backend/internal/config"

	ciokit "github.com/AddMile/backend/internal/kit/customerio"
)

type UserService interface {
	User(context.Context, user.UserParams) (user.User, error)
}

type EmailProvider interface {
	Identity(context.Context, ciokit.User) error
	TrackEvent(context.Context, ciokit.Event) error
}

type Processor struct {
	logger        *slog.Logger
	cfg           *config.Config
	emailProvider EmailProvider
	userService   UserService
}

func NewProcessor(
	logger *slog.Logger,
	cfg *config.Config,
	emailProvider EmailProvider,
	userService UserService,
) (*Processor, error) {
	return &Processor{
		logger:        logger,
		cfg:           cfg,
		emailProvider: emailProvider,
		userService:   userService,
	}, nil
}
