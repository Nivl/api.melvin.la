package mockdb

import db "github.com/Nivl/go-rest-tools/storage/db"
import mock "github.com/stretchr/testify/mock"

// Connection is an autogenerated mock type for the Connection type
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