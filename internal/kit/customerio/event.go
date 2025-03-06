package customerio

import (
	"fmt"
	"time"

	analytics "github.com/customerio/cdp-analytics-go"
)

type Event struct {
	MessageID    string
	AnonymousID  string
	UserID       string
	Event        string
	Timestamp    time.Time
	Properties   map[string]any
	Integrations map[string]any
}

func (e Event) Validate() error {
	if e.UserID == "" && e.AnonymousID == "" {
		return fmt.Errorf("%w: either user id or anonymous id must be set", ErrNoRequiredField)
	}

	if e.Event == "" {
		return fmt.Errorf("%w: event", ErrNoRequiredField)
	}

	return nil
}

func (e Event) toAnalyticsEvent() analytics.Track {
	return analytics.Track{
		MessageId:    e.MessageID,
		AnonymousId:  e.AnonymousID,
		UserId:       e.UserID,
		Event:        e.Event,
		Timestamp:    e.Timestamp,
		Properties:   e.Properties,
		Integrations: e.Integrations,
	}
}
