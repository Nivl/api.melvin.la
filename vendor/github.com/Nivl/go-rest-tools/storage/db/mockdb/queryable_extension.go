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
func (mdb *Queryable) ExpectGet(typ string, runnable func(args mock.Arguments)) *mock.Call {
	getCall := mdb.On("Get", mock.AnythingOfType(typ), StringType, StringType)
	getCall.Return(nil)
	if runnable != nil {
		getCall.Run(runnable)
	}
	return getCall
}

// ExpectGetNotFound is a helper that expects a not found on a Get
func (mdb *Queryable) ExpectGetNotFound(typ string) *mock.Call {
	getCall := mdb.On("Get", mock.AnythingOfType(typ), StringType, StringType)
	getCall.Return(sql.ErrNoRows)
	return getCall
}

// ExpectGetError is a helper that expects a connection error on a Get
func (mdb *Queryable) ExpectGetError(typ string) *mock.Call {
	getCall := mdb.On("Get", mock.AnythingOfType(typ), StringType, StringType)
	getCall.Return(serverError)
	return getCall
}

// ExpectSelect is an helper that expects a connection error on a Select
func (mdb *Queryable) ExpectSelect(typ string, runnable func(args mock.Arguments)) *mock.Call {
	selectCall := mdb.On("Select", mock.AnythingOfType(typ), StringType, InType, InType)
	selectCall.Return(nil)
	if runnable != nil {
		selectCall.Run(runnable)
	}
	return selectCall
}

// ExpectSelectError is an helper that expects a Select
func (mdb *Queryable) ExpectSelectError(typ string) *mock.Call {
	selectCall := mdb.On("Select", mock.AnythingOfType(typ), StringType, InType, InType)
	selectCall.Return(serverError)
	return selectCall
}

// ExpectDeletion is a helper that expects a deletion
func (mdb *Queryable) ExpectDeletion() *mock.Call {
	return mdb.On("Exec", StringType, StringType).Return(int64(1), nil)
}

// ExpectDeletionError is a helper that expects a deletion to fail
func (mdb *Queryable) ExpectDeletionError() *mock.Call {
	return mdb.On("Exec", StringType, StringType).Return(int64(0), serverError)
}

// ExpectInsert is a helper that expects an insertion
func (mdb *Queryable) ExpectInsert(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(int64(1), nil)
}

// ExpectInsertError is a helper that expects an insert to fail
func (mdb *Queryable) ExpectInsertError(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(int64(1), serverError)
}

// ExpectInsertConflict is a helper that expects a conflict on an insertion
func (mdb *Queryable) ExpectInsertConflict(typ string, fieldName string) *mock.Call {
	conflictError := newConflictError(fieldName)
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(int64(0), conflictError)
}

// ExpectUpdate is a helper that expects an update
func (mdb *Queryable) ExpectUpdate(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(int64(1), nil)
}

// ExpectUpdateConflict is a helper that expects a conflict on an update
func (mdb *Queryable) ExpectUpdateConflict(typ string, fieldName string) *mock.Call {
	conflictError := newConflictError(fieldName)
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(int64(0), conflictError)
}

// ExpectUpdateError is a helper that expects an update to fail
func (mdb *Queryable) ExpectUpdateError(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(int64(0), serverError)
}
