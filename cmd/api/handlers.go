package main

import (
	"encoding/json"
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

// sampleBook stands in for a database lookup until the service has a store.
func sampleBook(id int64) data.Book {
	return data.Book{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Echoes in the Darkness",
		Published: 2019,
		Pages:     300,
		Genres:    []string{"Fiction", "Thriller"},
		Rating:    4.5,
		Version:   1,
	}
}

func (app *application) listBooks(w http.ResponseWriter, _ *http.Request) {
	books := []data.Book{
		{ID: 1, CreatedAt: time.Now(), Title: "The Darkening of Tristram", Published: 1998, Pages: 300, Genres: []string{"Fiction", "Thriller"}, Rating: 4.5, Version: 1},
		{ID: 2, CreatedAt: time.Now(), Title: "The Legecy of Deckard Cain", Published: 2007, Pages: 432, Genres: []string{"Fiction", "Adventure"}, Rating: 4.9, Version: 1},
		{ID: 3, CreatedAt: time.Now(), Title: "The Black Soulstone", Version: 1},
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
	if err := app.writeJSON(w, http.StatusCreated, envelope{"book": input}); err != nil {
		app.serverError(w, err)
	}
}

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	id, err := bookID(r)
	if err != nil {
		badRequest(w)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"book": sampleBook(id)}); err != nil {
		app.serverError(w, err)
	}
}

func (app *application) updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := bookID(r)
	if err != nil {
		badRequest(w)
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

	book := sampleBook(id)
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
	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "deleted book " + strconv.FormatInt(id, 10)}); err != nil {
		app.serverError(w, err)
	}
}
