package mockdb

import (
	"errors"

	mock "github.com/stretchr/testify/mock"
)

// ExpectTransaction is a helper that expects a transaction
func (mdb *Connection) ExpectTransaction() (*Tx, *mock.Call) {
	tx := &Tx{}
	call := mdb.On("Beginx").Return(tx, nil)
	return tx, call
}

// ExpectTransactionError is a helper that expects a transaction to fail
func (mdb *Connection) ExpectTransactionError() *mock.Call {
	return mdb.On("Beginx").Return(nil, errors.New("cound not create transaction"))
}
