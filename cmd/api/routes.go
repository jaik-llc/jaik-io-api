package main

import "net/http"

// routes registers every endpoint this service publishes. The ingress sends all of
// api.jaik.io to this service, so a new endpoint only needs a line here.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", http.NotFound)
	mux.HandleFunc("GET /healthz", app.healthcheck)
	mux.HandleFunc("GET /healthcheck/db", app.dbHealthcheck)

	// api.jaik.io/books
	mux.HandleFunc("GET /books", app.listBooks)
	mux.HandleFunc("POST /books", app.createBook)
	mux.HandleFunc("GET /books/{id}", app.getBook)
	mux.HandleFunc("PUT /books/{id}", app.updateBook)
	mux.HandleFunc("DELETE /books/{id}", app.deleteBook)
	return mux
}
