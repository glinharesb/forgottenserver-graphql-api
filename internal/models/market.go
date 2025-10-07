package models

import (
	"context"
	"fmt"

	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
)

type MarketOffer struct {
	ID        int   `db:"id" json:"id"`
	PlayerID  int   `db:"player_id" json:"playerId"`
	Sale      bool  `db:"sale" json:"sale"`
	ItemType  int   `db:"itemtype" json:"itemType"`
	Amount    int   `db:"amount" json:"amount"`
	Created   int64 `db:"created" json:"created"`
	Anonymous bool  `db:"anonymous" json:"anonymous"`
	Price     int   `db:"price" json:"price"`
}

type MarketHistory struct {
	ID        int   `db:"id" json:"id"`
	PlayerID  int   `db:"player_id" json:"playerId"`
	Sale      bool  `db:"sale" json:"sale"`
	ItemType  int   `db:"itemtype" json:"itemType"`
	Amount    int   `db:"amount" json:"amount"`
	Price     int   `db:"price" json:"price"`
	ExpiresAt int64 `db:"expires_at" json:"expiresAt"`
	Inserted  int64 `db:"inserted" json:"inserted"`
	State     int   `db:"state" json:"state"`
}

type CreateMarketOfferInput struct {
	PlayerID  int
	Sale      bool
	ItemType  int
	Amount    int
	Price     int
	Anonymous bool
}

type MarketRepository struct {
	db *database.DB
}

func NewMarketRepository(db *database.DB) *MarketRepository {
	return &MarketRepository{db: db}
}

func (r *MarketRepository) GetOffers(ctx context.Context, itemType *int) ([]*MarketOffer, error) {
	var offers []*MarketOffer
	var query string

	if itemType != nil {
		query = `SELECT id, player_id, sale, itemtype, amount, created, anonymous, price
		         FROM market_offers WHERE itemtype = ? ORDER BY created DESC`
		if err := r.db.SelectContext(ctx, &offers, query, *itemType); err != nil {
			return nil, fmt.Errorf("failed to get market offers: %w", err)
		}
	} else {
		query = `SELECT id, player_id, sale, itemtype, amount, created, anonymous, price
		         FROM market_offers ORDER BY created DESC LIMIT 100`
		if err := r.db.SelectContext(ctx, &offers, query); err != nil {
			return nil, fmt.Errorf("failed to get market offers: %w", err)
		}
	}

	return offers, nil
}

func (r *MarketRepository) CreateOffer(ctx context.Context, input CreateMarketOfferInput) (*MarketOffer, error) {
	query := `INSERT INTO market_offers (player_id, sale, itemtype, amount, created, anonymous, price)
	          VALUES (?, ?, ?, ?, UNIX_TIMESTAMP(), ?, ?)`

	result, err := r.db.ExecContext(ctx, query, input.PlayerID, input.Sale, input.ItemType,
		input.Amount, input.Anonymous, input.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to create market offer: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	var offer MarketOffer
	query = `SELECT id, player_id, sale, itemtype, amount, created, anonymous, price
	         FROM market_offers WHERE id = ?`

	if err := r.db.GetContext(ctx, &offer, query, id); err != nil {
		return nil, fmt.Errorf("failed to get created offer: %w", err)
	}

	return &offer, nil
}

func (r *MarketRepository) GetHistory(ctx context.Context, playerID int) ([]*MarketHistory, error) {
	var history []*MarketHistory
	query := `SELECT id, player_id, sale, itemtype, amount, price, expires_at, inserted, state
	          FROM market_history WHERE player_id = ? ORDER BY inserted DESC`

	if err := r.db.SelectContext(ctx, &history, query, playerID); err != nil {
		return nil, fmt.Errorf("failed to get market history: %w", err)
	}

	return history, nil
}
