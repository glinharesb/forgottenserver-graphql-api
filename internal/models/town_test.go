package models

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTownRepository_GetByID(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTownRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "posx", "posy", "posz"}).
		AddRow(1, "Thais", 100, 200, 7)

	mock.ExpectQuery("SELECT id, name, posx, posy, posz FROM towns WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	town, err := repo.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 1, town.ID)
	assert.Equal(t, "Thais", town.Name)
	assert.Equal(t, 100, town.PosX)
	assert.Equal(t, 200, town.PosY)
	assert.Equal(t, 7, town.PosZ)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTownRepository_GetAll(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTownRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "posx", "posy", "posz"}).
		AddRow(1, "Thais", 100, 200, 7).
		AddRow(2, "Carlin", 150, 250, 7).
		AddRow(3, "Venore", 200, 300, 7)

	mock.ExpectQuery("SELECT id, name, posx, posy, posz FROM towns").
		WillReturnRows(rows)

	towns, err := repo.GetAll(context.Background())

	require.NoError(t, err)
	assert.Len(t, towns, 3)
	assert.Equal(t, "Thais", towns[0].Name)
	assert.Equal(t, "Carlin", towns[1].Name)
	assert.Equal(t, "Venore", towns[2].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTownRepository_Create(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTownRepository(db)

	input := CreateTownInput{
		Name: "Edron",
		PosX: 300,
		PosY: 400,
		PosZ: 7,
	}

	mock.ExpectExec("INSERT INTO towns").
		WithArgs(input.Name, input.PosX, input.PosY, input.PosZ).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id", "name", "posx", "posy", "posz"}).
		AddRow(1, input.Name, input.PosX, input.PosY, input.PosZ)

	mock.ExpectQuery("SELECT id, name, posx, posy, posz FROM towns WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	town, err := repo.Create(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, "Edron", town.Name)
	assert.Equal(t, 300, town.PosX)
	assert.Equal(t, 400, town.PosY)
	assert.Equal(t, 7, town.PosZ)
	assert.NoError(t, mock.ExpectationsWereMet())
}
