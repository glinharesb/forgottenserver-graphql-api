package models

import (
	"context"
	"fmt"

	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
)

type Town struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	PosX int    `db:"posx" json:"posX"`
	PosY int    `db:"posy" json:"posY"`
	PosZ int    `db:"posz" json:"posZ"`
}

type CreateTownInput struct {
	Name string
	PosX int
	PosY int
	PosZ int
}

type TownRepository struct {
	db *database.DB
}

func NewTownRepository(db *database.DB) *TownRepository {
	return &TownRepository{db: db}
}

func (r *TownRepository) GetByID(ctx context.Context, id int) (*Town, error) {
	var town Town
	query := `SELECT id, name, posx, posy, posz FROM towns WHERE id = ?`

	if err := r.db.GetContext(ctx, &town, query, id); err != nil {
		return nil, fmt.Errorf("failed to get town: %w", err)
	}

	return &town, nil
}

func (r *TownRepository) GetAll(ctx context.Context) ([]*Town, error) {
	var towns []*Town
	query := `SELECT id, name, posx, posy, posz FROM towns`

	if err := r.db.SelectContext(ctx, &towns, query); err != nil {
		return nil, fmt.Errorf("failed to get towns: %w", err)
	}

	return towns, nil
}

func (r *TownRepository) Create(ctx context.Context, input CreateTownInput) (*Town, error) {
	query := `INSERT INTO towns (name, posx, posy, posz) VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, input.Name, input.PosX, input.PosY, input.PosZ)
	if err != nil {
		return nil, fmt.Errorf("failed to create town: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(ctx, int(id))
}
