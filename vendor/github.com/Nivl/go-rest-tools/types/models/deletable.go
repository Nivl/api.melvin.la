package models

import "github.com/Nivl/go-rest-tools/storage/db"

// Deletable represents an objects that can be deleted from the database
type Deletable interface {
	Delete(q db.Queryable) error
}
