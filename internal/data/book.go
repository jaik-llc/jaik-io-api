// Package data holds the models returned by the API.
package data

import (
	"time"
)

// Book is a book as exposed by the /books endpoints.
type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Published int       `json:"published,omitempty"`
	Pages     int       `json:"pages,omitempty,string"` // change return data type to string
	Genres    []string  `json:"genres,omitempty"`
	Rating    float32   `json:"rating,omitempty"`
	Version   int32     `json:"-"`
}
