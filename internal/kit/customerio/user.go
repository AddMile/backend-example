package customerio

import (
	"fmt"
	"time"

	analytics "github.com/customerio/cdp-analytics-go"
)

type User struct {
	MessageID    string
	AnonymousID  string
	UserID       string
	Timestamp    time.Time
	Traits       map[string]any
	Integrations map[string]any
}

func (u User) Validate() error {
	if u.UserID == "" && u.AnonymousID == "" {
		return fmt.Errorf("%w: either user id or anonymous id must be set", ErrNoRequiredField)
	}

	return nil
}

func (u User) toUser() analytics.Identify {
	return analytics.Identify{
		MessageId:    u.MessageID,
		AnonymousId:  u.AnonymousID,
		UserId:       u.UserID,
		Timestamp:    u.Timestamp,
		Traits:       u.Traits,
		Integrations: u.Integrations,
	}
}
