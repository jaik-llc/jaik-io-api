package data

import (
	"context"
	"database/sql"
	"errors"
)

// BookModel wraps the connection pool for book queries.
type BookModel struct {
	DB *sql.DB
}

// Models groups every data model behind one handle that handlers use.
type Models struct {
	DB    *sql.DB
	Books BookModel
}

// NewModels builds the models on top of a database connection pool.
func NewModels(db *sql.DB) Models {
	return Models{
		DB:    db,
		Books: BookModel{DB: db},
	}
}

// Ping reports whether the database is reachable right now.
func (m Models) Ping(ctx context.Context) error {
	if m.DB == nil {
		return errors.New("no database configured")
	}
	return m.DB.PingContext(ctx)
}
