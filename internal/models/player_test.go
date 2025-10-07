package models

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlayerRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPlayerRepository(db)

	expectedPlayer := &Player{
		ID:         1,
		Name:       "TestPlayer",
		AccountID:  1,
		Level:      50,
		Vocation:   4,
		Health:     500,
		HealthMax:  500,
		Experience: 123456,
		MagLevel:   20,
		Mana:       300,
		ManaMax:    300,
		TownID:     1,
		Sex:        1,
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "group_id", "account_id", "level", "vocation", "health", "healthmax",
		"experience", "lookbody", "lookfeet", "lookhead", "looklegs", "looktype", "lookaddons",
		"maglevel", "mana", "manamax", "soul", "town_id", "posx", "posy", "posz", "cap", "sex",
		"lastlogin", "balance",
	}).AddRow(
		expectedPlayer.ID, expectedPlayer.Name, 1, expectedPlayer.AccountID, expectedPlayer.Level,
		expectedPlayer.Vocation, expectedPlayer.Health, expectedPlayer.HealthMax, expectedPlayer.Experience,
		0, 0, 0, 0, 136, 0, expectedPlayer.MagLevel, expectedPlayer.Mana, expectedPlayer.ManaMax,
		0, expectedPlayer.TownID, 0, 0, 0, 400, expectedPlayer.Sex, 0, 0,
	)

	mock.ExpectQuery("SELECT (.+) FROM players WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	player, err := repo.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, expectedPlayer.ID, player.ID)
	assert.Equal(t, expectedPlayer.Name, player.Name)
	assert.Equal(t, expectedPlayer.AccountID, player.AccountID)
	assert.Equal(t, expectedPlayer.Level, player.Level)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPlayerRepository_GetByID_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPlayerRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM players WHERE id = ?").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	player, err := repo.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, player)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPlayerRepository_GetByAccountID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPlayerRepository(db)

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

	players, err := repo.GetByAccountID(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, players, 2)
	assert.Equal(t, "Player1", players[0].Name)
	assert.Equal(t, "Player2", players[1].Name)
	assert.Equal(t, 1, players[0].AccountID)
	assert.Equal(t, 1, players[1].AccountID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPlayerRepository_GetByAccountID_NoPlayers(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPlayerRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "name", "group_id", "account_id", "level", "vocation", "health", "healthmax",
		"experience", "lookbody", "lookfeet", "lookhead", "looklegs", "looktype", "lookaddons",
		"maglevel", "mana", "manamax", "soul", "town_id", "posx", "posy", "posz", "cap", "sex",
		"lastlogin", "balance",
	})

	mock.ExpectQuery("SELECT (.+) FROM players WHERE account_id = ?").
		WithArgs(999).
		WillReturnRows(rows)

	players, err := repo.GetByAccountID(context.Background(), 999)

	require.NoError(t, err)
	assert.Len(t, players, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPlayerRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPlayerRepository(db)

	input := CreatePlayerInput{
		Name:      "NewPlayer",
		AccountID: 1,
		Sex:       1,
		Vocation:  4,
	}

	mock.ExpectExec("INSERT INTO players").
		WithArgs(input.Name, input.AccountID, input.Sex, input.Vocation).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock the GetByID call that happens after insert
	rows := sqlmock.NewRows([]string{
		"id", "name", "group_id", "account_id", "level", "vocation", "health", "healthmax",
		"experience", "lookbody", "lookfeet", "lookhead", "looklegs", "looktype", "lookaddons",
		"maglevel", "mana", "manamax", "soul", "town_id", "posx", "posy", "posz", "cap", "sex",
		"lastlogin", "balance",
	}).AddRow(
		1, input.Name, 1, input.AccountID, 1, input.Vocation, 150, 150, 0,
		0, 0, 0, 0, 136, 0, 0, 0, 0, 0, 1, 0, 0, 0, 400, input.Sex, 0, 0,
	)

	mock.ExpectQuery("SELECT (.+) FROM players WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	player, err := repo.Create(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, 1, player.ID)
	assert.Equal(t, input.Name, player.Name)
	assert.Equal(t, input.AccountID, player.AccountID)
	assert.Equal(t, input.Sex, player.Sex)
	assert.Equal(t, input.Vocation, player.Vocation)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPlayerRepository_Create_InvalidAccount(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	repo := NewPlayerRepository(db)

	input := CreatePlayerInput{
		Name:      "NewPlayer",
		AccountID: 999, // non-existent account
		Sex:       1,
		Vocation:  4,
	}

	// Simulate foreign key constraint violation
	mock.ExpectExec("INSERT INTO players").
		WithArgs(input.Name, input.AccountID, input.Sex, input.Vocation).
		WillReturnError(sql.ErrConnDone)

	player, err := repo.Create(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, player)
	assert.NoError(t, mock.ExpectationsWereMet())
}
