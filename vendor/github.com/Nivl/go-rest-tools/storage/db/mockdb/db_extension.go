package mockdb

import (
	"database/sql"
	"fmt"

	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/lib/pq"
	"github.com/stretchr/testify/mock"
)

var (
	// StringType represent a string argument
	StringType = mock.AnythingOfType("string")
)

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

// ExpectDeletion is a helper that expects a deletion
func (mdb *DB) ExpectDeletion() *mock.Call {
	return mdb.On("Exec", StringType, StringType).Return(nil, nil)
}

// ExpectInsert is a helper that expects an insertion
func (mdb *DB) ExpectInsert(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(nil, nil)
}

// ExpectInsertConflict is a helper that expects a conflict on an insertion
func (mdb *DB) ExpectInsertConflict(typ string, fieldName string) *mock.Call {
	conflictError := new(pq.Error)
	conflictError.Code = db.ErrDup
	conflictError.Message = "error: duplicate field"
	conflictError.Detail = fmt.Sprintf("Key (%s)=(Google) already exists.", fieldName)

	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(nil, conflictError)
}

// ExpectUpdate is a helper that expects an update
func (mdb *DB) ExpectUpdate(typ string) *mock.Call {
	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(nil, nil)
}

// ExpectUpdateConflict is a helper that expects a conflict on an update
func (mdb *DB) ExpectUpdateConflict(typ string, fieldName string) *mock.Call {
	conflictError := new(pq.Error)
	conflictError.Code = db.ErrDup
	conflictError.Message = "error: duplicate field"
	conflictError.Detail = fmt.Sprintf("Key (%s)=(Google) already exists.", fieldName)

	return mdb.On("NamedExec", StringType, mock.AnythingOfType(typ)).Return(nil, conflictError)
}
