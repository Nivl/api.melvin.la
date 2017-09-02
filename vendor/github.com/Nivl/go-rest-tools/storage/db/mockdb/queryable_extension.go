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
	return getCall.Once()
}

// ExpectGetID is a helper that expects a Get with a specific ID
func (mdb *Queryable) ExpectGetID(typ string, uuid string, runnable func(args mock.Arguments)) *mock.Call {
	getCall := mdb.On("Get", mock.AnythingOfType(typ), StringType, uuid)
	getCall.Return(nil)
	if runnable != nil {
		getCall.Run(runnable)
	}
	return getCall.Once()
}

// ExpectGetNotFound is a helper that expects a not found on a Get
func (mdb *Queryable) ExpectGetNotFound(typ string) *mock.Call {
	getCall := mdb.On("Get", mock.AnythingOfType(typ), StringType, StringType)
	getCall.Return(sql.ErrNoRows)
	return getCall.Once()
}

// ExpectGetIDNotFound is a helper that expects a not found on a Get with a specific ID
func (mdb *Queryable) ExpectGetIDNotFound(typ string, uuid string) *mock.Call {
	getCall := mdb.On("Get", mock.AnythingOfType(typ), StringType, uuid)
	getCall.Return(sql.ErrNoRows)
	return getCall.Once()
}

// ExpectGetIDError is a helper that expects a connection error on a Get with a specific ID
func (mdb *Queryable) ExpectGetIDError(typ string, uuid string) *mock.Call {
	getCall := mdb.On("Get", mock.AnythingOfType(typ), StringType, uuid)
	getCall.Return(serverError)
	return getCall.Once()
}

// ExpectGetNoParams is a helper that expects a Get with no params but the stmt
func (mdb *Queryable) ExpectGetNoParams(typ string, runnable func(args mock.Arguments)) *mock.Call {
	call := mdb.On("Get", mock.AnythingOfType(typ), StringType)
	call.Return(nil)
	if runnable != nil {
		call.Run(runnable)
	}
	return call.Once()
}

// ExpectGetNoParamsNotFound is a helper that expects a not found on a Get with no params but the stmt
func (mdb *Queryable) ExpectGetNoParamsNotFound(typ string) *mock.Call {
	call := mdb.On("Get", mock.AnythingOfType(typ), StringType)
	call.Return(sql.ErrNoRows)
	return call.Once()
}

// ExpectGetNoParamsError is a helper that expects a connection error on a Get with no params but the stmt
func (mdb *Queryable) ExpectGetNoParamsError(typ string) *mock.Call {
	call := mdb.On("Get", mock.AnythingOfType(typ), StringType)
	call.Return(serverError)
	return call.Once()
}

// ExpectGetError is a helper that expects a connection error on a Get
func (mdb *Queryable) ExpectGetError(typ string) *mock.Call {
	getCall := mdb.On("Get", mock.AnythingOfType(typ), StringType, StringType)
	getCall.Return(serverError)
	return getCall.Once()
}

// ExpectSelect is an helper that expects a connection error on a Select
func (mdb *Queryable) ExpectSelect(typ string, runnable func(args mock.Arguments)) *mock.Call {
	selectCall := mdb.On("Select", mock.AnythingOfType(typ), StringType, InType, InType)
	selectCall.Return(nil)
	if runnable != nil {
		selectCall.Run(runnable)
	}
	return selectCall.Once()
}

// ExpectSelectError is an helper that expects a Select
func (mdb *Queryable) ExpectSelectError(typ string) *mock.Call {
	selectCall := mdb.On("Select", mock.AnythingOfType(typ), StringType, InType, InType)
	selectCall.Return(serverError)
	return selectCall.Once()
}

// ExpectDeletion is a helper that expects a deletion
func (mdb *Queryable) ExpectDeletion() *mock.Call {
	return mdb.On("Exec", StringType, StringType).Return(int64(1), nil).Once()
}

// ExpectDeletionError is a helper that expects a deletion to fail
func (mdb *Queryable) ExpectDeletionError() *mock.Call {
	return mdb.On("Exec", StringType, StringType).Return(int64(0), serverError).Once()
}

// ExpectInsert is a helper that expects an insertion
func (mdb *Queryable) ExpectInsert(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(int64(1), nil).Once()
}

// ExpectInsertError is a helper that expects an insert to fail
func (mdb *Queryable) ExpectInsertError(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(int64(1), serverError).Once()
}

// ExpectInsertConflict is a helper that expects a conflict on an insertion
func (mdb *Queryable) ExpectInsertConflict(typ string, fieldName string) *mock.Call {
	conflictError := newConflictError(fieldName)
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(int64(0), conflictError).Once()
}

// ExpectUpdate is a helper that expects an update
func (mdb *Queryable) ExpectUpdate(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(int64(1), nil).Once()
}

// ExpectUpdateConflict is a helper that expects a conflict on an update
func (mdb *Queryable) ExpectUpdateConflict(typ string, fieldName string) *mock.Call {
	conflictError := newConflictError(fieldName)
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(int64(0), conflictError).Once()
}

// ExpectUpdateError is a helper that expects an update to fail
func (mdb *Queryable) ExpectUpdateError(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(int64(0), serverError).Once()
}
