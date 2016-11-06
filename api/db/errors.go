package db

import (
	"database/sql"

	"github.com/lib/pq"
)

const (
	// ErrDup contains the errcode of a unique constraint violation
	ErrDup = "23505"
)

// SQLIsDup check if an SQL error has been triggered by a duplicated data
func SQLIsDup(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == ErrDup
	}

	return false
}

// SQLNotFound checks if an error is triggered by an empty result
func SQLNotFound(err error) bool {
	return err != nil && err == sql.ErrNoRows
}
