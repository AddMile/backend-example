package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/AddMile/backend/internal/app/user"
	"github.com/AddMile/backend/internal/shared/event"

	ciokit "github.com/AddMile/backend/internal/kit/customerio"
)

func (p *Processor) EmailUserCreatedJob(ctx context.Context, event event.UserCreatedEvent) error {
	p.logger.Debug("user created event received", slog.Any("event", event))

	u, err := p.userService.User(ctx, user.UserParams{UserID: event.UserID})
	if err != nil {
		return err
	}

	cioUser := ciokit.User{
		UserID: u.ID.String(),
		Traits: map[string]any{
			"email": u.Email,
		},
	}

	err = p.emailProvider.Identity(ctx, cioUser)
	if err != nil {
		p.logger.Error("cannot create user in email provider", slog.Any("user", cioUser), slog.Any("error", err))

		return fmt.Errorf("cannot create user in email provider: %w", err)
	}

	cioEvent := ciokit.Event{
		Event:  "user_created",
		UserID: u.ID.String(),
	}

	err = p.emailProvider.TrackEvent(ctx, cioEvent)
	if err != nil {
		p.logger.Error("cannot send user created", slog.Any("event", cioEvent), slog.Any("error", err))

		return fmt.Errorf("cannot send user created: %w", err)
	}

	return nil
}
