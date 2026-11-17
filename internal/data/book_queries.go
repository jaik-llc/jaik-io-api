package data

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var ErrRecordNotFound = errors.New("record not found")

func (m BookModel) Insert(book *Book) error {
	query := `
        INSERT INTO books (title, published, pages, genres, rating)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, version`

	args := []any{
		book.Title,
		book.Published,
		book.Pages,
		pq.Array(book.Genres),
		book.Rating,
	}
	return m.DB.QueryRow(query, args...).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Version,
	)
}

func (m BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, title, published, pages,
               genres, version, rating
        FROM books
        WHERE id = $1`

	var book Book

	err := m.DB.QueryRow(query, id).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Title,
		&book.Published,
		&book.Pages,
		pq.Array(&book.Genres),
		&book.Version,
		&book.Rating,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &book, nil
}

func (m BookModel) Update(book *Book) error {
	query := `
        UPDATE books
        SET title = $1,
            published = $2,
            pages = $3,
            genres = $4,
            rating = $5,
            version = version + 1
        WHERE id = $6
        RETURNING version`

	args := []any{
		book.Title,
		book.Published,
		book.Pages,
		pq.Array(book.Genres),
		book.Rating,
		book.ID,
	}

	err := m.DB.QueryRow(query, args...).Scan(&book.Version)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrRecordNotFound
	}

	return err
}

func (m BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM books WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m BookModel) GetAll() ([]*Book, error) {
	query := `
        SELECT id, created_at, title, published, pages,
               genres, version, rating
        FROM books
        ORDER BY id`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []*Book{}

	for rows.Next() {
		var book Book

		err := rows.Scan(
			&book.ID,
			&book.CreatedAt,
			&book.Title,
			&book.Published,
			&book.Pages,
			pq.Array(&book.Genres),
			&book.Version,
			&book.Rating,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
