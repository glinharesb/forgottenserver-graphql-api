package models

import (
	"context"
	"fmt"

	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
)

type PlayerDeath struct {
	PlayerID            int    `db:"player_id" json:"playerId"`
	Time                int64  `db:"time" json:"time"`
	Level               int    `db:"level" json:"level"`
	KilledBy            string `db:"killed_by" json:"killedBy"`
	IsPlayer            bool   `db:"is_player" json:"isPlayer"`
	MostDamageBy        string `db:"mostdamage_by" json:"mostDamageBy"`
	MostDamageIsPlayer  bool   `db:"mostdamage_is_player" json:"mostDamageIsPlayer"`
}

type PlayerDeathRepository struct {
	db *database.DB
}

func NewPlayerDeathRepository(db *database.DB) *PlayerDeathRepository {
	return &PlayerDeathRepository{db: db}
}

func (r *PlayerDeathRepository) GetByPlayerID(ctx context.Context, playerID int) ([]*PlayerDeath, error) {
	var deaths []*PlayerDeath
	query := `SELECT player_id, time, level, killed_by, is_player, mostdamage_by, mostdamage_is_player
	          FROM player_deaths WHERE player_id = ? ORDER BY time DESC`

	if err := r.db.SelectContext(ctx, &deaths, query, playerID); err != nil {
		return nil, fmt.Errorf("failed to get player deaths: %w", err)
	}

	return deaths, nil
}
