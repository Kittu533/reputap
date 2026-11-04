package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func Connect(dsn string) (*DB, error) {
	sdb, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sdb.SetMaxOpenConns(25)
	sdb.SetMaxIdleConns(5)
	sdb.SetConnMaxLifetime(5 * time.Minute)
	sdb.SetConnMaxIdleTime(2 * time.Minute)

	if err := sdb.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("postgresql database connected successfully")
	return &DB{DB: sdb}, nil
}

func (db *DB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.DB.GetContext(ctx, dest, db.Rebind(query), args...)
}

func (db *DB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return db.DB.SelectContext(ctx, dest, db.Rebind(query), args...)
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.DB.ExecContext(ctx, db.Rebind(query), args...)
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRowContext(ctx, db.Rebind(query), args...)
}

func (db *DB) InsertID(ctx context.Context, query string, args ...interface{}) (uint64, error) {
	var id uint64
	q := strings.TrimSpace(query)
	if !strings.HasSuffix(strings.ToUpper(q), "RETURNING ID") {
		q += " RETURNING id"
	}
	err := db.DB.QueryRowContext(ctx, db.Rebind(q), args...).Scan(&id)
	return id, err
}
