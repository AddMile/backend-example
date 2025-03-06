package user

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/AddMile/backend/internal/shared"

	pgkit "github.com/AddMile/backend/internal/kit/pg"
)

type Storage struct {
	db *pgkit.DB
}

func NewStorage(
	db *pgkit.DB,
) *Storage {
	return &Storage{
		db: db,
	}
}

type userRecord struct {
	ID        uuid.UUID       `db:"id"`
	Email     string          `db:"email"`
	Language  shared.Language `db:"language"`
	CreatedAt time.Time       `db:"created_at"`
}

func (ur userRecord) toDomain() User {
	return User(ur)
}

const upsertUserQuery = `
  INSERT INTO users
  (email, language)
  VALUES (@email, @language)
  ON CONFLICT (email)
  DO UPDATE SET
    language   = EXCLUDED.language
  RETURNING id;
`

func (s *Storage) UpsertUser(ctx context.Context, data UpsertUserData) (uuid.UUID, error) {
	var userID uuid.UUID

	args := pgkit.Args{
		"email":    data.Email,
		"language": data.Language,
	}

	if err := pgkit.QueryRow(ctx, s.db, upsertUserQuery, &userID, args); err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

const userQuery = `
  SELECT
    id,
    email,
    language,
    created_at
  FROM users
  WHERE id = @id;
`

func (s *Storage) User(ctx context.Context, data UserData) (User, error) {
	args := pgkit.Args{"id": data.UserID}

	record, err := pgkit.QueryRowStruct[userRecord](ctx, s.db, userQuery, args)
	if err != nil {
		return User{}, err
	}

	return record.toDomain(), nil
}
