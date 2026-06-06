package db

import (
	"database/sql"
	"telemetry-mend-server/internal/models"
	"context"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func InitDB(ctx context.Context, dsn string) (*bun.DB, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())

	if err := models.CreateTables(ctx, db); err != nil {
		return nil, err
	}

	return db, nil
}
