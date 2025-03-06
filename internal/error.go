package internal

import (
	"errors"
	"fmt"

	pgkit "github.com/AddMile/backend/internal/kit/pg"
)

var (
	// client errors.
	ErrAlreadyExists     = fmt.Errorf("already exists")
	ErrAlreadyInUse      = fmt.Errorf("already in use")
	ErrNotFound          = fmt.Errorf("not found")
	ErrValidation        = fmt.Errorf("validation failed")
	ErrConditionViolated = fmt.Errorf("conditiion violated")

	// dependency errors.
	ErrNilLogger      = errors.New("nil logger")
	ErrNilConfig      = errors.New("nil config")
	ErrNilStorage     = errors.New("nil storage")
	ErrNilGeoStorage  = errors.New("nil geo storage")
	ErrNilQueue       = errors.New("nil queue")
	ErrNilValidator   = errors.New("nil validator")
	ErrNilUserService = errors.New("nil user service")
)

func ClientErr(err error) bool {
	switch {
	case errors.Is(err, ErrAlreadyExists),
		errors.Is(err, ErrAlreadyInUse),
		errors.Is(err, ErrNotFound),
		errors.Is(err, ErrValidation),
		errors.Is(err, ErrConditionViolated),
		errors.Is(err, pgkit.ErrNoRows):

		return true
	}

	return false
}
