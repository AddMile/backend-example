package google

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

type TokenConfig struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

type Token struct {
	*oauth2.Token
	Refreshed bool
}

func RefreshTokenIfNeeded(ctx context.Context, c *oauth2.Config, t TokenConfig) (Token, error) {
	token := &oauth2.Token{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
	}

	if token.Valid() {
		return Token{Token: token}, nil
	}

	ts := c.TokenSource(ctx, &oauth2.Token{RefreshToken: t.RefreshToken})
	newToken, err := ts.Token()
	if err != nil {
		return Token{}, fmt.Errorf("failed to refresh token: %w", err)
	}

	return Token{Token: newToken, Refreshed: true}, nil
}
