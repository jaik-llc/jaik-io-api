// Package data holds the models returned by the API.
package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

// Book is a book as exposed by the /books endpoints.
type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Published int       `json:"published,omitempty"`
	Pages     int       `json:"pages,omitempty,string"`
	Genres    []string  `json:"genres,omitempty"`
	Rating    float32   `json:"rating,omitempty"`
	Version   int32     `json:"-"`
}

const bookColumns = `id, created_at, title, published, pages, genres, rating, version`

func scanBook(row interface{ Scan(...any) error }) (Book, error) {
	var book Book
	var genres pq.StringArray
	err := row.Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Title,
		&book.Published,
		&book.Pages,
		&genres,
		&book.Rating,
		&book.Version,
	)
	if err != nil {
		return Book{}, err
	}
	book.Genres = []string(genres)
	return book, nil
}

func (m BookModel) GetAll(ctx context.Context) ([]Book, error) {
	rows, err := m.DB.QueryContext(ctx, `SELECT `+bookColumns+` FROM books ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]Book, 0)
	for rows.Next() {
		book, err := scanBook(rows)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (m BookModel) Get(ctx context.Context, id int64) (Book, error) {
	row := m.DB.QueryRowContext(ctx, `SELECT `+bookColumns+` FROM books WHERE id = $1`, id)
	return scanBook(row)
}

func (m BookModel) Insert(ctx context.Context, book Book) (Book, error) {
	genres := book.Genres
	if genres == nil {
		genres = []string{}
	}
	row := m.DB.QueryRowContext(ctx, `
		INSERT INTO books (title, published, pages, genres, rating)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+bookColumns,
		book.Title, book.Published, book.Pages, pq.Array(genres), book.Rating,
	)
	return scanBook(row)
}

func (m BookModel) Update(ctx context.Context, book Book) (Book, error) {
	row := m.DB.QueryRowContext(ctx, `
		UPDATE books
		SET title = $2, published = $3, pages = $4, genres = $5, rating = $6, version = version + 1
		WHERE id = $1
		RETURNING `+bookColumns,
		book.ID, book.Title, book.Published, book.Pages, pq.Array(book.Genres), book.Rating,
	)
	return scanBook(row)
}

func (m BookModel) Delete(ctx context.Context, id int64) error {
	result, err := m.DB.ExecContext(ctx, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
