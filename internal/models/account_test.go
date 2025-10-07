package models

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockDB(t *testing.T) (*database.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	db := &database.DB{DB: sqlxDB}

	return db, mock
}

func TestAccountRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewAccountRepository(db)

	expectedAccount := &Account{
		ID:            1,
		Name:          "testuser",
		Password:      "hashedpassword",
		Email:         "test@example.com",
		Type:          1,
		PremiumEndsAt: 0,
		Creation:      1234567890,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "password", "secret", "type", "premium_ends_at", "email", "creation"}).
		AddRow(expectedAccount.ID, expectedAccount.Name, expectedAccount.Password, nil, expectedAccount.Type, expectedAccount.PremiumEndsAt, expectedAccount.Email, expectedAccount.Creation)

	mock.ExpectQuery("SELECT id, name, password, secret, type, premium_ends_at, email, creation FROM accounts WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	account, err := repo.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, expectedAccount.ID, account.ID)
	assert.Equal(t, expectedAccount.Name, account.Name)
	assert.Equal(t, expectedAccount.Email, account.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountRepository_GetByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewAccountRepository(db)

	mock.ExpectQuery("SELECT id, name, password, secret, type, premium_ends_at, email, creation FROM accounts WHERE id = ?").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	account, err := repo.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, account)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountRepository_GetAll(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewAccountRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "password", "secret", "type", "premium_ends_at", "email", "creation"}).
		AddRow(1, "user1", "pass1", nil, 1, 0, "user1@example.com", 1234567890).
		AddRow(2, "user2", "pass2", nil, 1, 0, "user2@example.com", 1234567891)

	mock.ExpectQuery("SELECT id, name, password, secret, type, premium_ends_at, email, creation FROM accounts LIMIT ?").
		WithArgs(10).
		WillReturnRows(rows)

	accounts, err := repo.GetAll(context.Background(), 10)

	require.NoError(t, err)
	assert.Len(t, accounts, 2)
	assert.Equal(t, "user1", accounts[0].Name)
	assert.Equal(t, "user2", accounts[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewAccountRepository(db)

	input := CreateAccountInput{
		Name:     "newuser",
		Password: "password123",
		Email:    "newuser@example.com",
	}

	mock.ExpectExec("INSERT INTO accounts").
		WithArgs(input.Name, input.Password, input.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock the GetByID call that happens after insert
	rows := sqlmock.NewRows([]string{"id", "name", "password", "secret", "type", "premium_ends_at", "email", "creation"}).
		AddRow(1, input.Name, input.Password, nil, 1, 0, input.Email, 1234567890)

	mock.ExpectQuery("SELECT id, name, password, secret, type, premium_ends_at, email, creation FROM accounts WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	account, err := repo.Create(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, 1, account.ID)
	assert.Equal(t, input.Name, account.Name)
	assert.Equal(t, input.Email, account.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountRepository_Create_DuplicateName(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewAccountRepository(db)

	input := CreateAccountInput{
		Name:     "existinguser",
		Password: "password123",
		Email:    "test@example.com",
	}

	mock.ExpectExec("INSERT INTO accounts").
		WithArgs(input.Name, input.Password, input.Email).
		WillReturnError(sql.ErrConnDone) // Simulate duplicate key error

	account, err := repo.Create(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, account)
	assert.NoError(t, mock.ExpectationsWereMet())
}
