package pg

import "github.com/jackc/pgx/v5"

type Args map[string]any

func (a Args) toPGx() pgx.NamedArgs {
	if a == nil {
		return pgx.NamedArgs{}
	}

	result := make(pgx.NamedArgs, len(a))
	for k, v := range a {
		result[k] = v
	}

	return result
}
