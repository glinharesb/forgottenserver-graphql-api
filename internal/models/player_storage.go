package models

import (
	"context"
	"fmt"

	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
)

type PlayerStorage struct {
	PlayerID int `db:"player_id" json:"playerId"`
	Key      int `db:"key" json:"key"`
	Value    int `db:"value" json:"value"`
}

type PlayerStorageRepository struct {
	db *database.DB
}

func NewPlayerStorageRepository(db *database.DB) *PlayerStorageRepository {
	return &PlayerStorageRepository{db: db}
}

func (r *PlayerStorageRepository) GetByPlayerID(ctx context.Context, playerID int) ([]*PlayerStorage, error) {
	var storage []*PlayerStorage
	query := `SELECT player_id, key, value FROM player_storage WHERE player_id = ?`

	if err := r.db.SelectContext(ctx, &storage, query, playerID); err != nil {
		return nil, fmt.Errorf("failed to get player storage: %w", err)
	}

	return storage, nil
}
