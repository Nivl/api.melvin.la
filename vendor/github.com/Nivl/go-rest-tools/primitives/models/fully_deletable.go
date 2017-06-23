package models

import "github.com/jmoiron/sqlx"

// Deletable represents an objects that can be deleted from the database
type Deletable interface {
	Delete(q *sqlx.DB) error
}
