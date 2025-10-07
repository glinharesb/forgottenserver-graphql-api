package models

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHouseRepository_GetByID(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewHouseRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "owner", "paid", "warnings", "name", "rent", "town_id", "bid", "bid_end",
		"last_bid", "highest_bidder", "size", "beds",
	}).AddRow(1, 0, 0, 0, "Test House", 1000, 1, 0, 0, 0, 0, 100, 2)

	mock.ExpectQuery("SELECT (.+) FROM houses WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(rows)

	house, err := repo.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 1, house.ID)
	assert.Equal(t, "Test House", house.Name)
	assert.Equal(t, 1000, house.Rent)
	assert.Equal(t, 1, house.TownID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHouseRepository_GetByTownID(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewHouseRepository(db)

	t.Run("GetAllHouses", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "owner", "paid", "warnings", "name", "rent", "town_id", "bid", "bid_end",
			"last_bid", "highest_bidder", "size", "beds",
		}).
			AddRow(1, 0, 0, 0, "House 1", 1000, 1, 0, 0, 0, 0, 100, 2).
			AddRow(2, 0, 0, 0, "House 2", 1500, 2, 0, 0, 0, 0, 150, 3)

		mock.ExpectQuery("SELECT (.+) FROM houses$").
			WillReturnRows(rows)

		houses, err := repo.GetByTownID(context.Background(), nil)

		require.NoError(t, err)
		assert.Len(t, houses, 2)
		assert.Equal(t, "House 1", houses[0].Name)
		assert.Equal(t, "House 2", houses[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetHousesByTown", func(t *testing.T) {
		townID := 1
		rows := sqlmock.NewRows([]string{
			"id", "owner", "paid", "warnings", "name", "rent", "town_id", "bid", "bid_end",
			"last_bid", "highest_bidder", "size", "beds",
		}).AddRow(1, 0, 0, 0, "House 1", 1000, 1, 0, 0, 0, 0, 100, 2)

		mock.ExpectQuery("SELECT (.+) FROM houses WHERE town_id = \\?").
			WithArgs(townID).
			WillReturnRows(rows)

		houses, err := repo.GetByTownID(context.Background(), &townID)

		require.NoError(t, err)
		assert.Len(t, houses, 1)
		assert.Equal(t, 1, houses[0].TownID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestHouseRepository_PlaceBid(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewHouseRepository(db)

	mock.ExpectExec("UPDATE houses SET bid = \\?, highest_bidder = \\?, last_bid = UNIX_TIMESTAMP\\(\\) WHERE id = \\?").
		WithArgs(5000, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rows := sqlmock.NewRows([]string{
		"id", "owner", "paid", "warnings", "name", "rent", "town_id", "bid", "bid_end",
		"last_bid", "highest_bidder", "size", "beds",
	}).AddRow(1, 0, 0, 0, "Test House", 1000, 1, 5000, 0, 1234567890, 1, 100, 2)

	mock.ExpectQuery("SELECT (.+) FROM houses WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(rows)

	house, err := repo.PlaceBid(context.Background(), 1, 1, 5000)

	require.NoError(t, err)
	assert.Equal(t, 5000, house.Bid)
	assert.Equal(t, 1, house.HighestBidder)
	assert.NoError(t, mock.ExpectationsWereMet())
}
