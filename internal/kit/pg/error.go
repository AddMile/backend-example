package pg

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNoRows = errors.New("no rows found")

func RecordNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func RecordAlreadyExists(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code)
}
