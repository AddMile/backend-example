package event

import "github.com/google/uuid"

type UserCreatedEvent struct {
	UserID uuid.UUID `json:"user_id"`
}
