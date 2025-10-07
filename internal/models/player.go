package models

import (
	"context"
	"fmt"

	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
)

type Player struct {
	ID         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	GroupID    int    `db:"group_id" json:"groupId"`
	AccountID  int    `db:"account_id" json:"accountId"`
	Level      int    `db:"level" json:"level"`
	Vocation   int    `db:"vocation" json:"vocation"`
	Health     int    `db:"health" json:"health"`
	HealthMax  int    `db:"healthmax" json:"healthMax"`
	Experience int64  `db:"experience" json:"experience"`
	LookBody   int    `db:"lookbody" json:"lookBody"`
	LookFeet   int    `db:"lookfeet" json:"lookFeet"`
	LookHead   int    `db:"lookhead" json:"lookHead"`
	LookLegs   int    `db:"looklegs" json:"lookLegs"`
	LookType   int    `db:"looktype" json:"lookType"`
	LookAddons int    `db:"lookaddons" json:"lookAddons"`
	MagLevel   int    `db:"maglevel" json:"magLevel"`
	Mana       int    `db:"mana" json:"mana"`
	ManaMax    int    `db:"manamax" json:"manaMax"`
	Soul       int    `db:"soul" json:"soul"`
	TownID     int    `db:"town_id" json:"townId"`
	PosX       int    `db:"posx" json:"posX"`
	PosY       int    `db:"posy" json:"posY"`
	PosZ       int    `db:"posz" json:"posZ"`
	Cap        int    `db:"cap" json:"cap"`
	Sex        int    `db:"sex" json:"sex"`
	LastLogin  int64  `db:"lastlogin" json:"lastLogin"`
	Balance    int64  `db:"balance" json:"balance"`
}

type CreatePlayerInput struct {
	Name      string
	AccountID int
	Sex       int
	Vocation  int
}

type PlayerRepository struct {
	db *database.DB
}

func NewPlayerRepository(db *database.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) GetByID(ctx context.Context, id int) (*Player, error) {
	var player Player
	query := `
		SELECT id, name, group_id, account_id, level, vocation, health, healthmax,
		       experience, lookbody, lookfeet, lookhead, looklegs, looktype, lookaddons,
		       maglevel, mana, manamax, soul, town_id, posx, posy, posz, cap, sex,
		       lastlogin, balance
		FROM players
		WHERE id = ?
	`

	if err := r.db.GetContext(ctx, &player, query, id); err != nil {
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	return &player, nil
}

func (r *PlayerRepository) GetByAccountID(ctx context.Context, accountID int) ([]*Player, error) {
	var players []*Player
	query := `
		SELECT id, name, group_id, account_id, level, vocation, health, healthmax,
		       experience, lookbody, lookfeet, lookhead, looklegs, looktype, lookaddons,
		       maglevel, mana, manamax, soul, town_id, posx, posy, posz, cap, sex,
		       lastlogin, balance
		FROM players
		WHERE account_id = ?
	`

	if err := r.db.SelectContext(ctx, &players, query, accountID); err != nil {
		return nil, fmt.Errorf("failed to get players: %w", err)
	}

	return players, nil
}

func (r *PlayerRepository) Create(ctx context.Context, input CreatePlayerInput) (*Player, error) {
	query := `
		INSERT INTO players (name, account_id, sex, vocation)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, input.Name, input.AccountID, input.Sex, input.Vocation)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(ctx, int(id))
}
