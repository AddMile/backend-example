package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/AddMile/backend/internal/shared"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Language  shared.Language
	CreatedAt time.Time
}
