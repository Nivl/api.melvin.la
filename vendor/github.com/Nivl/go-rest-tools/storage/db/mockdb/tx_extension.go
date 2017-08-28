package mockdb

import (
	"errors"

	mock "github.com/stretchr/testify/mock"
)

// ExpectCommit is a helper that expects a transaction to Commit
func (mdb *Tx) ExpectCommit() *mock.Call {
	return mdb.On("Commit").Return(nil)
}

// ExpectCommitError is a helper that expects a commit to fail
func (mdb *Tx) ExpectCommitError() *mock.Call {
	return mdb.On("Commit").Return(errors.New("could not commit"))
}

// ExpectRollback is a helper that expects a transaction to Rollback
func (mdb *Tx) ExpectRollback() *mock.Call {
	return mdb.On("Rollback").Return(nil)
}

// ExpectRollbackError is a helper that expects a Rollback to fail
func (mdb *Tx) ExpectRollbackError() *mock.Call {
	return mdb.On("Rollback").Return(errors.New("could not Rollback"))
}
