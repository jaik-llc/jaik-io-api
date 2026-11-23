package data

import "database/sql"

type Models struct {
	Book BookModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Book: BookModel{DB: db},
	}
}