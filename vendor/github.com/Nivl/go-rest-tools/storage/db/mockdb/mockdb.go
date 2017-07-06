package mockdb

import (
	"database/sql"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/jmoiron/sqlx"
)

// MockDB represents a mocked DB
type MockDB struct {
	Mock  sqlmock.Sqlmock
	SQLDB *sql.DB
	DB    db.DB
}

// New creates a new MockDB
func New() (*MockDB, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	return &MockDB{
		Mock:  mock,
		SQLDB: db,
		DB:    sqlx.NewDb(db, "sqlmock"),
	}, nil
}
