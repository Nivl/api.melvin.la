package db

import "github.com/jmoiron/sqlx"

// Setup setup the database connection and init the Writer
func Setup(uri string) error {
	db, err := sqlx.Connect("postgres", uri)
	if err != nil {
		return err
	}

	// Unsafe returns a version of DB which will silently succeed to scan when
	// columns in the SQL result have no fields in the destination struct.
	Writer = db.Unsafe()
	return nil
}
