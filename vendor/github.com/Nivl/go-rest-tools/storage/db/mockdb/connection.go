package mockdb

import db "github.com/Nivl/go-rest-tools/storage/db"
import mock "github.com/stretchr/testify/mock"
import sql "database/sql"

// Connection is a mock type for the Connection type
// NOT auto generated - do not add the Queryable method here as it will
// overwite the ones defined in the Queryabl file.
type Connection struct {
	Queryable
	mock.Mock
}

// Beginx provides a mock function with given fields:
func (_m *Connection) Beginx() (db.Tx, error) {
	ret := _m.Called()

	var r0 db.Tx
	if rf, ok := ret.Get(0).(func() db.Tx); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Tx)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *Connection) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DSN provides a mock function with given fields:
func (_m *Connection) DSN() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SQL provides a mock function with given fields:
func (_m *Connection) SQL() *sql.DB {
	ret := _m.Called()

	var r0 *sql.DB
	if rf, ok := ret.Get(0).(func() *sql.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.DB)
		}
	}

	return r0
}
