package db

import (
	"database/sql"

	"github.com/lib/pq"
)

const (
	// ErrDup contains the errcode of a unique constraint violation
	ErrDup = "23505"
)

// IsDup check if an SQL error has been triggered by a duplicated data
func IsDup(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == ErrDup
	}

	return false
}

// IsNotFound checks if an error is triggered by an empty result
func IsNotFound(err error) bool {
	return err != nil && err == sql.ErrNoRows
}
