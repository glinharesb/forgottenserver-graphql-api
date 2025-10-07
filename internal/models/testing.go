package models

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
	"github.com/jmoiron/sqlx"
)

// NewMockDB creates a new mock database for testing
func NewMockDB() (*database.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	db := &database.DB{DB: sqlxDB}

	return db, mock, nil
}

// CloseMockDB closes the mock database
func CloseMockDB(db *database.DB) error {
	if db != nil && db.DB != nil {
		return db.Close()
	}
	return nil
}

// Helper type to simplify cleanup in tests
type MockDB struct {
	DB   *database.DB
	Mock sqlmock.Sqlmock
}

// NewMock creates a new MockDB with proper cleanup
func NewMock() (*MockDB, error) {
	var mockSQL *sql.DB
	var mock sqlmock.Sqlmock
	var err error

	mockSQL, mock, err = sqlmock.New()
	if err != nil {
		return nil, err
	}

	sqlxDB := sqlx.NewDb(mockSQL, "sqlmock")
	db := &database.DB{DB: sqlxDB}

	return &MockDB{
		DB:   db,
		Mock: mock,
	}, nil
}

// Close closes the mock database
func (m *MockDB) Close() error {
	if m.DB != nil {
		return m.DB.Close()
	}
	return nil
}
