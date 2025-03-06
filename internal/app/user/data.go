package user

import (
	"github.com/google/uuid"

	"github.com/AddMile/backend/internal/shared"
)

type UpsertUserData struct {
	Email    string
	Language shared.Language
}

type UserData struct {
	UserID uuid.UUID
}
