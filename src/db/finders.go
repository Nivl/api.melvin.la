package db

// Get is the same as sqlx.Get() but do not returns an error on empty results
func Get(dest interface{}, query string, args ...interface{}) error {
	err := Con().Get(dest, query, args...)

	if SQLNotFound(err) {
		return nil
	}

	return err
}
