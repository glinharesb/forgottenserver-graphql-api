package graph

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
	"github.com/glinharesb/forgottenserver-graphql-api/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestResolver(t *testing.T) (*Resolver, sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	db := &database.DB{DB: sqlxDB}

	resolver := NewResolver(db)

	cleanup := func() {
		db.Close()
	}

	return resolver, mock, cleanup
}

// Query Resolver Tests

func TestQueryResolver_Account(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "password", "secret", "type", "premium_ends_at", "email", "creation"}).
		AddRow(1, "testuser", "hashedpass", nil, 1, 0, "test@example.com", 1234567890)

	mock.ExpectQuery("SELECT id, name, password, secret, type, premium_ends_at, email, creation FROM accounts WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	account, err := resolver.Query().Account(context.Background(), "1")

	require.NoError(t, err)
	assert.Equal(t, 1, account.ID)
	assert.Equal(t, "testuser", account.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryResolver_Account_InvalidID(t *testing.T) {
	resolver, _, cleanup := setupTestResolver(t)
	defer cleanup()

	account, err := resolver.Query().Account(context.Background(), "invalid")

	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Contains(t, err.Error(), "invalid account id")
}

func TestQueryResolver_Accounts(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "password", "secret", "type", "premium_ends_at", "email", "creation"}).
		AddRow(1, "user1", "pass1", nil, 1, 0, "user1@example.com", 1234567890).
		AddRow(2, "user2", "pass2", nil, 1, 0, "user2@example.com", 1234567891)

	limit := 10
	mock.ExpectQuery("SELECT id, name, password, secret, type, premium_ends_at, email, creation FROM accounts LIMIT ?").
		WithArgs(limit).
		WillReturnRows(rows)

	accounts, err := resolver.Query().Accounts(context.Background(), &limit)

	require.NoError(t, err)
	assert.Len(t, accounts, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryResolver_Player(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "name", "group_id", "account_id", "level", "vocation", "health", "healthmax",
		"experience", "lookbody", "lookfeet", "lookhead", "looklegs", "looktype", "lookaddons",
		"maglevel", "mana", "manamax", "soul", "town_id", "posx", "posy", "posz", "cap", "sex",
		"lastlogin", "balance",
	}).AddRow(1, "TestPlayer", 1, 1, 50, 4, 500, 500, 123456, 0, 0, 0, 0, 136, 0, 20, 300, 300, 0, 1, 0, 0, 0, 400, 1, 0, 0)

	mock.ExpectQuery("SELECT (.+) FROM players WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	player, err := resolver.Query().Player(context.Background(), "1")

	require.NoError(t, err)
	assert.Equal(t, 1, player.ID)
	assert.Equal(t, "TestPlayer", player.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryResolver_Players(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "name", "group_id", "account_id", "level", "vocation", "health", "healthmax",
		"experience", "lookbody", "lookfeet", "lookhead", "looklegs", "looktype", "lookaddons",
		"maglevel", "mana", "manamax", "soul", "town_id", "posx", "posy", "posz", "cap", "sex",
		"lastlogin", "balance",
	}).
		AddRow(1, "Player1", 1, 1, 20, 1, 200, 200, 5000, 0, 0, 0, 0, 136, 0, 5, 50, 50, 0, 1, 0, 0, 0, 400, 1, 0, 0).
		AddRow(2, "Player2", 1, 1, 30, 2, 300, 300, 10000, 0, 0, 0, 0, 136, 0, 10, 100, 100, 0, 1, 0, 0, 0, 400, 0, 0, 0)

	mock.ExpectQuery("SELECT (.+) FROM players WHERE account_id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	players, err := resolver.Query().Players(context.Background(), "1")

	require.NoError(t, err)
	assert.Len(t, players, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Mutation Resolver Tests

func TestMutationResolver_CreateAccount(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	input := models.CreateAccountInput{
		Name:     "newuser",
		Password: "password123",
		Email:    "newuser@example.com",
	}

	mock.ExpectExec("INSERT INTO accounts").
		WithArgs(input.Name, input.Password, input.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id", "name", "password", "secret", "type", "premium_ends_at", "email", "creation"}).
		AddRow(1, input.Name, input.Password, nil, 1, 0, input.Email, 1234567890)

	mock.ExpectQuery("SELECT id, name, password, secret, type, premium_ends_at, email, creation FROM accounts WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	account, err := resolver.Mutation().CreateAccount(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, 1, account.ID)
	assert.Equal(t, input.Name, account.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMutationResolver_CreatePlayer(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	input := models.CreatePlayerInput{
		Name:      "NewPlayer",
		AccountID: 1,
		Sex:       1,
		Vocation:  4,
	}

	mock.ExpectExec("INSERT INTO players").
		WithArgs(input.Name, input.AccountID, input.Sex, input.Vocation).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{
		"id", "name", "group_id", "account_id", "level", "vocation", "health", "healthmax",
		"experience", "lookbody", "lookfeet", "lookhead", "looklegs", "looktype", "lookaddons",
		"maglevel", "mana", "manamax", "soul", "town_id", "posx", "posy", "posz", "cap", "sex",
		"lastlogin", "balance",
	}).AddRow(1, input.Name, 1, input.AccountID, 1, input.Vocation, 150, 150, 0, 0, 0, 0, 0, 136, 0, 0, 0, 0, 0, 1, 0, 0, 0, 400, input.Sex, 0, 0)

	mock.ExpectQuery("SELECT (.+) FROM players WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	player, err := resolver.Mutation().CreatePlayer(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, 1, player.ID)
	assert.Equal(t, input.Name, player.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Field Resolver Tests

func TestAccountResolver_Players(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	account := &models.Account{
		ID:   1,
		Name: "testuser",
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "group_id", "account_id", "level", "vocation", "health", "healthmax",
		"experience", "lookbody", "lookfeet", "lookhead", "looklegs", "looktype", "lookaddons",
		"maglevel", "mana", "manamax", "soul", "town_id", "posx", "posy", "posz", "cap", "sex",
		"lastlogin", "balance",
	}).AddRow(1, "Player1", 1, 1, 20, 1, 200, 200, 5000, 0, 0, 0, 0, 136, 0, 5, 50, 50, 0, 1, 0, 0, 0, 400, 1, 0, 0)

	mock.ExpectQuery("SELECT (.+) FROM players WHERE account_id = ?").
		WithArgs(account.ID).
		WillReturnRows(rows)

	players, err := resolver.Account().Players(context.Background(), account)

	require.NoError(t, err)
	assert.Len(t, players, 1)
	assert.Equal(t, "Player1", players[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPlayerResolver_Account(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	player := &models.Player{
		ID:        1,
		Name:      "TestPlayer",
		AccountID: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "password", "secret", "type", "premium_ends_at", "email", "creation"}).
		AddRow(1, "testuser", "hashedpass", nil, 1, 0, "test@example.com", 1234567890)

	mock.ExpectQuery("SELECT id, name, password, secret, type, premium_ends_at, email, creation FROM accounts WHERE id = ?").
		WithArgs(player.AccountID).
		WillReturnRows(rows)

	account, err := resolver.Player().Account(context.Background(), player)

	require.NoError(t, err)
	assert.Equal(t, 1, account.ID)
	assert.Equal(t, "testuser", account.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPlayerResolver_Account_NotFound(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	player := &models.Player{
		ID:        1,
		Name:      "TestPlayer",
		AccountID: 999,
	}

	mock.ExpectQuery("SELECT id, name, password, secret, type, premium_ends_at, email, creation FROM accounts WHERE id = ?").
		WithArgs(player.AccountID).
		WillReturnError(sql.ErrNoRows)

	account, err := resolver.Player().Account(context.Background(), player)

	assert.Error(t, err)
	assert.Nil(t, account)
	assert.NoError(t, mock.ExpectationsWereMet())
}
