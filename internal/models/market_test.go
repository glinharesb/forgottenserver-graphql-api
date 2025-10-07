package models

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarketRepository_GetOffers(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMarketRepository(db)

	t.Run("GetAllOffers", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "player_id", "sale", "itemtype", "amount", "created", "anonymous", "price",
		}).
			AddRow(1, 1, true, 2160, 10, 1234567890, false, 1000).
			AddRow(2, 2, false, 2160, 5, 1234567891, true, 1100)

		mock.ExpectQuery("SELECT (.+) FROM market_offers ORDER BY created DESC LIMIT 100").
			WillReturnRows(rows)

		offers, err := repo.GetOffers(context.Background(), nil)

		require.NoError(t, err)
		assert.Len(t, offers, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetOffersByItemType", func(t *testing.T) {
		itemType := 2160
		rows := sqlmock.NewRows([]string{
			"id", "player_id", "sale", "itemtype", "amount", "created", "anonymous", "price",
		}).AddRow(1, 1, true, 2160, 10, 1234567890, false, 1000)

		mock.ExpectQuery("SELECT (.+) FROM market_offers WHERE itemtype = \\? ORDER BY created DESC").
			WithArgs(itemType).
			WillReturnRows(rows)

		offers, err := repo.GetOffers(context.Background(), &itemType)

		require.NoError(t, err)
		assert.Len(t, offers, 1)
		assert.Equal(t, 2160, offers[0].ItemType)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMarketRepository_CreateOffer(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMarketRepository(db)

	input := CreateMarketOfferInput{
		PlayerID:  1,
		Sale:      true,
		ItemType:  2160,
		Amount:    10,
		Price:     1000,
		Anonymous: false,
	}

	mock.ExpectExec("INSERT INTO market_offers").
		WithArgs(input.PlayerID, input.Sale, input.ItemType, input.Amount, input.Anonymous, input.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{
		"id", "player_id", "sale", "itemtype", "amount", "created", "anonymous", "price",
	}).AddRow(1, input.PlayerID, input.Sale, input.ItemType, input.Amount, 1234567890, input.Anonymous, input.Price)

	mock.ExpectQuery("SELECT (.+) FROM market_offers WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(rows)

	offer, err := repo.CreateOffer(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, 1, offer.ID)
	assert.Equal(t, 2160, offer.ItemType)
	assert.Equal(t, 1000, offer.Price)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarketRepository_GetHistory(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMarketRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "player_id", "sale", "itemtype", "amount", "price", "expires_at", "inserted", "state",
	}).
		AddRow(1, 1, true, 2160, 10, 1000, 1234567890, 1234567800, 1).
		AddRow(2, 1, false, 2152, 5, 500, 1234567891, 1234567801, 2)

	mock.ExpectQuery("SELECT (.+) FROM market_history WHERE player_id = \\? ORDER BY inserted DESC").
		WithArgs(1).
		WillReturnRows(rows)

	history, err := repo.GetHistory(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, history, 2)
	assert.Equal(t, 1, history[0].PlayerID)
	assert.Equal(t, 2160, history[0].ItemType)
	assert.NoError(t, mock.ExpectationsWereMet())
}
