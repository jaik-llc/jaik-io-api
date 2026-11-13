package main

import (
	"io"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestApp() *application {
	return &application{config: config{env: "test"}, logger: log.New(io.Discard, "", 0)}
}

func TestRoutes(t *testing.T) {
	h := newTestApp().routes()
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
