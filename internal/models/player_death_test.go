package models

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlayerDeathRepository_GetByPlayerID(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPlayerDeathRepository(db)

	rows := sqlmock.NewRows([]string{
		"player_id", "time", "level", "killed_by", "is_player", "mostdamage_by", "mostdamage_is_player",
	}).
		AddRow(1, 1234567890, 50, "Dragon", false, "Dragon", false).
		AddRow(1, 1234567880, 49, "Player2", true, "Player2", true)

	mock.ExpectQuery("SELECT (.+) FROM player_deaths WHERE player_id = \\? ORDER BY time DESC").
		WithArgs(1).
		WillReturnRows(rows)

	deaths, err := repo.GetByPlayerID(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, deaths, 2)
	assert.Equal(t, "Dragon", deaths[0].KilledBy)
	assert.False(t, deaths[0].IsPlayer)
	assert.Equal(t, "Player2", deaths[1].KilledBy)
	assert.True(t, deaths[1].IsPlayer)
	assert.NoError(t, mock.ExpectationsWereMet())
}
