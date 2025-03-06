package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

var ErrDSNRequired = errors.New("dsn required")

func New(dsn string) (*DB, error) {
	if dsn == "" {
		return nil, ErrDSNRequired
	}

	ctx := context.Background()

	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// by default pgx postgres timestamptz into local
	// with this hook we ensure timestamptz to time.Time is always in UTC.
	conf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterType(&pgtype.Type{
			Name:  "timestamptz",
			OID:   pgtype.TimestamptzOID,
			Codec: &pgtype.TimestamptzCodec{ScanLocation: time.UTC},
		})

		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("creating pgpool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("pinging db: %w", err)
	}

	return &DB{Pool: pool}, nil
}
