package slack

import (
	"errors"
	"fmt"
)

var ErrNoRequiredField = errors.New("no required field")

type Event struct {
	Endpoint string
	Message  string
}

func (e Event) Validate() error {
	if e.Endpoint == "" || e.Message == "" {
		return fmt.Errorf("%w: endpoint and message must be set", ErrNoRequiredField)
	}

	return nil
}
