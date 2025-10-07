package models

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountBanRepository_GetByAccountID(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewAccountBanRepository(db)

	rows := sqlmock.NewRows([]string{"account_id", "reason", "banned_at", "expires_at", "banned_by"}).
		AddRow(1, "Botting", 1234567890, 1234567900, 1).
		AddRow(1, "Hacking", 1234567800, 1234567850, 2)

	mock.ExpectQuery("SELECT account_id, reason, banned_at, expires_at, banned_by FROM account_bans WHERE account_id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	bans, err := repo.GetByAccountID(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, bans, 2)
	assert.Equal(t, "Botting", bans[0].Reason)
	assert.Equal(t, "Hacking", bans[1].Reason)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountBanRepository_Create(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewAccountBanRepository(db)

	input := BanAccountInput{
		AccountID: 1,
		Reason:    "Botting",
		ExpiresAt: 1234567900,
		BannedBy:  1,
	}

	mock.ExpectExec("INSERT INTO account_bans").
		WithArgs(input.AccountID, input.Reason, input.ExpiresAt, input.BannedBy).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rows := sqlmock.NewRows([]string{"account_id", "reason", "banned_at", "expires_at", "banned_by"}).
		AddRow(input.AccountID, input.Reason, 1234567890, input.ExpiresAt, input.BannedBy)

	mock.ExpectQuery("SELECT account_id, reason, banned_at, expires_at, banned_by FROM account_bans WHERE account_id = ?").
		WithArgs(input.AccountID).
		WillReturnRows(rows)

	ban, err := repo.Create(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, 1, ban.AccountID)
	assert.Equal(t, "Botting", ban.Reason)
	assert.Equal(t, int64(1234567900), ban.ExpiresAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}
