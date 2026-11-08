package main

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jaik-llc/jaik-io-api/internal/data"
)

// fakeDriver stands in for PostgreSQL: it connects when ok is true and refuses otherwise.
type fakeDriver struct{ ok bool }

func (d fakeDriver) Open(string) (driver.Conn, error) {
	if !d.ok {
		return nil, errors.New("connection refused")
	}
	return fakeConn{}, nil
}

type fakeConn struct{}

func (fakeConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("unsupported") }
func (fakeConn) Close() error                        { return nil }
func (fakeConn) Begin() (driver.Tx, error)           { return nil, errors.New("unsupported") }

func init() {
	sql.Register("fake-up", fakeDriver{ok: true})
	sql.Register("fake-down", fakeDriver{ok: false})
}

func TestDBHealthcheck(t *testing.T) {
	tests := []struct {
		name     string
		models   func(t *testing.T) data.Models
		wantCode int
		wantBody string
	}{
		{"database up", func(t *testing.T) data.Models {
			db, err := sql.Open("fake-up", "")
			if err != nil {
				t.Fatal(err)
			}
			return data.NewModels(db)
		}, 200, "true"},
		{"database down", func(t *testing.T) data.Models {
			db, err := sql.Open("fake-down", "")
			if err != nil {
				t.Fatal(err)
			}
			return data.NewModels(db)
		}, 401, "false"},
		{"no database configured", func(*testing.T) data.Models { return data.Models{} }, 401, "false"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := newTestApp()
			app.models = tc.models(t)
			rec := httptest.NewRecorder()
			app.routes().ServeHTTP(rec, httptest.NewRequest("GET", "/healthcheck/db", nil))
			if rec.Code != tc.wantCode {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantCode)
			}
			if got := strings.TrimSpace(rec.Body.String()); got != tc.wantBody {
				t.Errorf("body = %q, want %q", got, tc.wantBody)
			}
		})
	}
}
