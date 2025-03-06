package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func QueryRow[T any](ctx context.Context, db *DB, q string, dest *T, args Args) error {
	err := db.QueryRow(ctx, q, args.toPGx()).Scan(dest)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRows
		}

		return fmt.Errorf("failed to execute QueryRow for query %q with args %v: %w", q, args, err)
	}

	return nil
}

func QueryRows[T any](ctx context.Context, db *DB, q string, args Args) ([]T, error) {
	rows, err := db.Query(ctx, q, args.toPGx())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collected, err := pgx.CollectRows(rows, pgx.RowTo[T])
	if err != nil {
		return nil, fmt.Errorf("failed to execute QueryRows for query %q with args %v: %w", q, args, err)
	}

	return collected, nil
}

func QueryRowStruct[T any](ctx context.Context, db *DB, q string, args Args) (T, error) {
	var record T

	rows, err := db.Query(ctx, q, args.toPGx())
	if err != nil {
		return record, err
	}
	defer rows.Close()

	collected, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return record, ErrNoRows
		}

		return record, fmt.Errorf("failed to execute QueryRowStruct for query %q with args %v: %w", q, args, err)
	}

	return collected, nil
}

func QueryRowsStruct[T any](ctx context.Context, db *DB, q string, args Args) ([]T, error) {
	rows, err := db.Query(ctx, q, args.toPGx())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collected, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRows
		}

		return nil, fmt.Errorf("failed to execute QueryRowsStruct for query %q with args %v: %w", q, args, err)
	}

	if len(collected) == 0 {
		return nil, nil
	}

	return collected, nil
}

func Exec(ctx context.Context, db *DB, q string, args Args) error {
	_, err := db.Exec(ctx, q, args.toPGx())
	if err != nil {
		return fmt.Errorf("failed to execute Exec %q with args %v: %w", q, args, err)
	}

	return nil
}

func Batch(ctx context.Context, db *DB, q string, args []Args) error {
	var batch pgx.Batch
	for _, i := range args {
		batch.Queue(q, i.toPGx())
	}

	err := db.SendBatch(ctx, &batch).Close()
	if err != nil {
		return fmt.Errorf("failed to execute Batch %q with args %v: %w", q, args, err)
	}

	return nil
}

func WithTx(ctx context.Context, db *DB, fn func(tx pgx.Tx) error) error {
	err := pgx.BeginFunc(ctx, db, fn)
	if err != nil {
		return fmt.Errorf("failed to execute WithTx: %w", err)
	}

	return nil
}

func QueryRowTx[T any](ctx context.Context, tx pgx.Tx, q string, dest *T, args Args) error {
	err := tx.QueryRow(ctx, q, args.toPGx()).Scan(dest)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRows
		}

		return fmt.Errorf("failed to execute QueryRowTx for query %q with args %v: %w", q, args, err)
	}

	return nil
}

func QueryRowsTx[T any](ctx context.Context, tx pgx.Tx, q string, args Args) ([]T, error) {
	rows, err := tx.Query(ctx, q, args.toPGx())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collected, err := pgx.CollectRows(rows, pgx.RowTo[T])
	if err != nil {
		return nil, fmt.Errorf("failed to execute QueryRowsTx for query %q with args %v: %w", q, args, err)
	}

	return collected, nil
}

func QueryRowStructTx[T any](ctx context.Context, tx pgx.Tx, q string, args Args) (T, error) {
	var record T

	rows, err := tx.Query(ctx, q, args.toPGx())
	if err != nil {
		return record, err
	}
	defer rows.Close()

	collected, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return record, ErrNoRows
		}

		return record, fmt.Errorf("failed to execute QueryRowStructTx for query %q with args %v: %w", q, args, err)
	}

	return collected, nil
}

func QueryRowsStructTx[T any](ctx context.Context, tx pgx.Tx, q string, args Args) ([]T, error) {
	rows, err := tx.Query(ctx, q, args.toPGx())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collected, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRows
		}

		return nil, fmt.Errorf("failed to execute QueryRowsStructTx for query %q with args %v: %w", q, args, err)
	}

	if len(collected) == 0 {
		return nil, nil
	}

	return collected, nil
}

func ExecTx(ctx context.Context, tx pgx.Tx, q string, args Args) error {
	_, err := tx.Exec(ctx, q, args.toPGx())
	if err != nil {
		return fmt.Errorf("failed to execute execTx %q with args %v: %w", q, args, err)
	}

	return nil
}

func BatchTx(ctx context.Context, tx pgx.Tx, q string, args []Args) error {
	var batch pgx.Batch
	for _, i := range args {
		batch.Queue(q, i.toPGx())
	}

	err := tx.SendBatch(ctx, &batch).Close()
	if err != nil {
		return fmt.Errorf("failed to execute batchTx %q with args %v: %w", q, args, err)
	}

	return nil
}
