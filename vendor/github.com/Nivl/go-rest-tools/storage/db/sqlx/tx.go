package sqlx

import (
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/jmoiron/sqlx"
)

var _ db.Tx = (*Tx)(nil)

// NewTx create a new transaction
func NewTx(con *sqlx.DB) (*Tx, error) {
	tx, err := con.Beginx()
	if err != nil {
		return nil, err
	}
	return &Tx{
		Queryable: NewQueryable(tx),
		tx:        tx,
	}, nil
}

// Tx implements the db.Tx interface for sqlx
type Tx struct {
	*Queryable
	tx *sqlx.Tx
}

// Commit commits the transaction
func (con *Tx) Commit() error {
	return con.tx.Commit()
}

// Rollback roll backs the transaction
func (con *Tx) Rollback() error {
	return con.tx.Rollback()
}
