package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jaik-llc/jaik-io-api/internal/data"
)

func newTestApp(t *testing.T) *application {
	db, err := sql.Open("fake-books", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	app := &application{
		config: config{env: "test"},
		logger: log.New(io.Discard, "", 0),
		models: data.NewModels(db),
	}
	return app
}

type booksTestDriver struct{}

func (booksTestDriver) Open(string) (driver.Conn, error) { return booksTestConn{}, nil }

type booksTestConn struct{}

func (booksTestConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("unsupported") }
func (booksTestConn) Close() error                        { return nil }
func (booksTestConn) Begin() (driver.Tx, error)           { return nil, errors.New("unsupported") }
func (booksTestConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return &booksTestRows{
		columns: []string{"id", "created_at", "title", "published", "pages", "genres", "rating", "version"},
		values:  []driver.Value{int64(7), time.Now(), "Test book", int64(2024), int64(250), "{Fiction}", float64(4.5), int64(1)},
	}, nil
}
func (booksTestConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	return driver.RowsAffected(1), nil
}

type booksTestRows struct {
	columns []string
	values  []driver.Value
	done    bool
}

func (r *booksTestRows) Columns() []string { return r.columns }
func (r *booksTestRows) Close() error      { return nil }
func (r *booksTestRows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}
	copy(dest, r.values)
	r.done = true
	return nil
}

func init() {
	sql.Register("fake-books", booksTestDriver{})
}

func TestRoutes(t *testing.T) {
	app := newTestApp(t)
	h := app.routes()
	tests := []struct {
		method, path, body string
		want               int
	}{
		{"GET", "/healthz", "", 200},
		{"GET", "/books", "", 200},
		{"GET", "/books/7", "", 200},
		{"GET", "/books/abc", "", 400},
		{"POST", "/books", `{"title":"x"}`, 201},
		{"POST", "/books", `{"nope":1}`, 400},
		{"PUT", "/books/7", `{"title":"y"}`, 200},
		{"DELETE", "/books/7", "", 200},
		{"GET", "/v1/books", "", 404},
	}
	for _, tc := range tests {
		req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != tc.want {
			t.Errorf("%s %s = %d, want %d", tc.method, tc.path, rec.Code, tc.want)
		}
	}
}
