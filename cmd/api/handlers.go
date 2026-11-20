// Command api serves the jaik.io HTTP API (api.jaik.io).
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jaik-llc/jaik-io-api/internal/data"
)

func (app *application) healthcheck(w http.ResponseWriter, _ *http.Request) {
	js, err := json.Marshal(map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(append(js, '\n'))
}

// dbHealthcheck verifies the database connection. It answers 200 with the JSON body true when the
// database responds, and 401 with false when it cannot be reached. It is deliberately not used by
// the Kubernetes probes: a database outage must not restart the pods.
func (app *application) dbHealthcheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	status, ok := http.StatusOK, true
	if err := app.models.Ping(ctx); err != nil {
		app.logger.Printf("database healthcheck failed: %v", err)
		status, ok = http.StatusUnauthorized, false
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = fmt.Fprintln(w, ok)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	app.logger.Printf("error: %v", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func bookID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func (app *application) listBooks(w http.ResponseWriter, _ *http.Request) {
	books, err := app.models.Books.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"books": books}); err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createBook(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title     string   `json:"title"`
		Published int      `json:"published"`
		Pages     int      `json:"pages"`
		Genres    []string `json:"genres"`
		Rating    float32  `json:"rating"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		badRequest(w)
		return
	}

	book := &data.Book{
		Title:     input.Title,
		Published: input.Published,
		Pages:     input.Pages,
		Genres:    input.Genres,
		Rating:    input.Rating,
	}

	if err := app.models.Books.Insert(book); err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/books/%d", book.ID))
	if err := app.writeJSON(w, http.StatusCreated, envelope{"book": book}); err != nil {
		app.serverError(w, err)
	}
}

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	id, err := bookID(r)
	if err != nil {
		badRequest(w)
		return
	}

	book, err := app.models.Books.Get(id)
	if err != nil {
		if err.Error() == "record not found" {
			http.NotFound(w, r)
			return
		}
		app.serverError(w, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"book": book}); err != nil {
		app.serverError(w, err)
	}
}

func (app *application) updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := bookID(r)
	if err != nil {
		badRequest(w)
		return
	}

	book, err := app.models.Books.Get(id)
	if err != nil {
		if err.Error() == "record not found" {
			http.NotFound(w, r)
			return
		}
		app.serverError(w, err)
		return
	}

	var input struct {
		Title     *string  `json:"title"`
		Published *int     `json:"published"`
		Pages     *int     `json:"pages"`
		Genres    []string `json:"genres"`
		Rating    *float32 `json:"rating"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		badRequest(w)
		return
	}

	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Published != nil {
		book.Published = *input.Published
	}
	if input.Pages != nil {
		book.Pages = *input.Pages
	}
	if len(input.Genres) > 0 {
		book.Genres = input.Genres
	}
	if input.Rating != nil {
		book.Rating = *input.Rating
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"book": book}); err != nil {
		app.serverError(w, err)
	}
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := bookID(r)
	if err != nil {
		badRequest(w)
		return
	}
	if err := app.models.Books.Delete(id); err != nil {
		if err.Error() == "record not found" {
			http.NotFound(w, r)
			return
		}
		app.serverError(w, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "deleted book " + strconv.FormatInt(id, 10)}); err != nil {
		app.serverError(w, err)
	}
}
