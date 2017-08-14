package mockdb

import (
	"database/sql"
	"fmt"

	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/lib/pq"
	"github.com/stretchr/testify/mock"
)

var (
	// StringType represents a string argument
	StringType = mock.AnythingOfType("string")
	// InType represents an int argument
	InType = mock.AnythingOfType("int")

	// serverError represents a database connection error
	serverError = &pq.Error{
		Code:    "08006",
		Message: "error: connection failure",
		Detail:  "the connection to the database failed",
	}
)

func newConflictError(fieldName string) *pq.Error {
	return &pq.Error{
		Code:    db.ErrDup,
		Message: "error: duplicate field",
		Detail:  fmt.Sprintf("Key (%s)=(Google) already exists.", fieldName),
	}
}

// ExpectGet is a helper that expects a Get
func (mdb *DB) ExpectGet(typ string, runnable func(args mock.Arguments)) *mock.Call {
	getCall := mdb.On("Get", mock.AnythingOfType(typ), StringType, StringType)
	getCall.Return(nil)
	if runnable != nil {
		getCall.Run(runnable)
	}
	return getCall
}

// ExpectGetNotFound is a helper that expects a not found on a Get
func (mdb *DB) ExpectGetNotFound(typ string) *mock.Call {
	getCall := mdb.On("Get", mock.AnythingOfType(typ), StringType, StringType)
	getCall.Return(sql.ErrNoRows)
	return getCall
}

// ExpectGetError is a helper that expects a connection error on a Get
func (mdb *DB) ExpectGetError(typ string) *mock.Call {
	getCall := mdb.On("Get", mock.AnythingOfType(typ), StringType, StringType)
	getCall.Return(serverError)
	return getCall
}

// ExpectSelectError is an helper that expects a connection error on a Select
func (mdb *DB) ExpectSelectError(typ string) *mock.Call {
	getCall := mdb.On("Select", mock.AnythingOfType(typ), StringType, InType, InType)
	getCall.Return(serverError)
	return getCall
}

// ExpectDeletion is a helper that expects a deletion
func (mdb *DB) ExpectDeletion() *mock.Call {
	return mdb.On("Exec", StringType, StringType).Return(nil, nil)
}

// ExpectDeletionError is a helper that expects a deletion to fail
func (mdb *DB) ExpectDeletionError() *mock.Call {
	return mdb.On("Exec", StringType, StringType).Return(nil, serverError)
}

// ExpectInsert is a helper that expects an insertion
func (mdb *DB) ExpectInsert(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(nil, nil)
}

// ExpectInsertError is a helper that expects an insert to fail
func (mdb *DB) ExpectInsertError(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(nil, serverError)
}

// ExpectInsertConflict is a helper that expects a conflict on an insertion
func (mdb *DB) ExpectInsertConflict(typ string, fieldName string) *mock.Call {
	conflictError := newConflictError(fieldName)
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(nil, conflictError)
}

// ExpectUpdate is a helper that expects an update
func (mdb *DB) ExpectUpdate(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(nil, nil)
}

// ExpectUpdateConflict is a helper that expects a conflict on an update
func (mdb *DB) ExpectUpdateConflict(typ string, fieldName string) *mock.Call {
	conflictError := newConflictError(fieldName)
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(nil, conflictError)
}

// ExpectUpdateError is a helper that expects an update to fail
func (mdb *DB) ExpectUpdateError(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(nil, serverError)
}
