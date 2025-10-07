package models

import (
	"context"
	"fmt"

	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
)

type AccountStorage struct {
	AccountID int `db:"account_id" json:"accountId"`
	Key       int `db:"key" json:"key"`
	Value     int `db:"value" json:"value"`
}

type VipEntry struct {
	AccountID   int    `db:"account_id" json:"accountId"`
	PlayerID    int    `db:"player_id" json:"playerId"`
	Description string `db:"description" json:"description"`
	Icon        int    `db:"icon" json:"icon"`
	Notify      bool   `db:"notify" json:"notify"`
}

type AccountStorageRepository struct {
	db *database.DB
}

func NewAccountStorageRepository(db *database.DB) *AccountStorageRepository {
	return &AccountStorageRepository{db: db}
}

func (r *AccountStorageRepository) GetByAccountID(ctx context.Context, accountID int) ([]*AccountStorage, error) {
	var storage []*AccountStorage
	query := `SELECT account_id, key, value FROM account_storage WHERE account_id = ?`

	if err := r.db.SelectContext(ctx, &storage, query, accountID); err != nil {
		return nil, fmt.Errorf("failed to get account storage: %w", err)
	}

	return storage, nil
}

func (r *AccountStorageRepository) GetVipList(ctx context.Context, accountID int) ([]*VipEntry, error) {
	var vipList []*VipEntry
	query := `SELECT account_id, player_id, description, icon, notify
	          FROM account_viplist WHERE account_id = ?`

	if err := r.db.SelectContext(ctx, &vipList, query, accountID); err != nil {
		return nil, fmt.Errorf("failed to get vip list: %w", err)
	}

	return vipList, nil
}
