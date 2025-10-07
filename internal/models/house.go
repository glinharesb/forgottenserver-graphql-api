package models

import (
	"context"
	"fmt"

	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
)

type House struct {
	ID             int    `db:"id" json:"id"`
	Owner          int    `db:"owner" json:"owner"`
	Paid           int    `db:"paid" json:"paid"`
	Warnings       int    `db:"warnings" json:"warnings"`
	Name           string `db:"name" json:"name"`
	Rent           int    `db:"rent" json:"rent"`
	TownID         int    `db:"town_id" json:"townId"`
	Bid            int    `db:"bid" json:"bid"`
	BidEnd         int    `db:"bid_end" json:"bidEnd"`
	LastBid        int    `db:"last_bid" json:"lastBid"`
	HighestBidder  int    `db:"highest_bidder" json:"highestBidder"`
	Size           int    `db:"size" json:"size"`
	Beds           int    `db:"beds" json:"beds"`
}

type HouseList struct {
	HouseID int    `db:"house_id" json:"houseId"`
	ListID  int    `db:"listid" json:"listId"`
	List    string `db:"list" json:"list"`
}

type HouseRepository struct {
	db *database.DB
}

func NewHouseRepository(db *database.DB) *HouseRepository {
	return &HouseRepository{db: db}
}

func (r *HouseRepository) GetByID(ctx context.Context, id int) (*House, error) {
	var house House
	query := `SELECT id, owner, paid, warnings, name, rent, town_id, bid, bid_end, last_bid,
	          highest_bidder, size, beds FROM houses WHERE id = ?`

	if err := r.db.GetContext(ctx, &house, query, id); err != nil {
		return nil, fmt.Errorf("failed to get house: %w", err)
	}

	return &house, nil
}

func (r *HouseRepository) GetByTownID(ctx context.Context, townID *int) ([]*House, error) {
	var houses []*House
	var query string

	if townID != nil {
		query = `SELECT id, owner, paid, warnings, name, rent, town_id, bid, bid_end, last_bid,
		         highest_bidder, size, beds FROM houses WHERE town_id = ?`
		if err := r.db.SelectContext(ctx, &houses, query, *townID); err != nil {
			return nil, fmt.Errorf("failed to get houses: %w", err)
		}
	} else {
		query = `SELECT id, owner, paid, warnings, name, rent, town_id, bid, bid_end, last_bid,
		         highest_bidder, size, beds FROM houses`
		if err := r.db.SelectContext(ctx, &houses, query); err != nil {
			return nil, fmt.Errorf("failed to get houses: %w", err)
		}
	}

	return houses, nil
}

func (r *HouseRepository) PlaceBid(ctx context.Context, houseID, playerID, bidAmount int) (*House, error) {
	query := `UPDATE houses SET bid = ?, highest_bidder = ?, last_bid = UNIX_TIMESTAMP()
	          WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, bidAmount, playerID, houseID)
	if err != nil {
		return nil, fmt.Errorf("failed to place bid: %w", err)
	}

	return r.GetByID(ctx, houseID)
}
