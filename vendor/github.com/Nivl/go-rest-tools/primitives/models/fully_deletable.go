package models

// FullyDeletable represents an objects that can be deleted from the database
type FullyDeletable interface {
	FullyDelete() error
}
