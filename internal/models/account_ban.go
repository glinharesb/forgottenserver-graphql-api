package models

import (
	"context"
	"fmt"

	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
)

type AccountBan struct {
	AccountID int    `db:"account_id" json:"accountId"`
	Reason    string `db:"reason" json:"reason"`
	BannedAt  int64  `db:"banned_at" json:"bannedAt"`
	ExpiresAt int64  `db:"expires_at" json:"expiresAt"`
	BannedBy  int    `db:"banned_by" json:"bannedBy"`
}

type BanAccountInput struct {
	AccountID int
	Reason    string
	ExpiresAt int64
	BannedBy  int
}

type AccountBanRepository struct {
	db *database.DB
}

func NewAccountBanRepository(db *database.DB) *AccountBanRepository {
	return &AccountBanRepository{db: db}
}

func (r *AccountBanRepository) GetByAccountID(ctx context.Context, accountID int) ([]*AccountBan, error) {
	var bans []*AccountBan
	query := `SELECT account_id, reason, banned_at, expires_at, banned_by FROM account_bans WHERE account_id = ?`

	if err := r.db.SelectContext(ctx, &bans, query, accountID); err != nil {
		return nil, fmt.Errorf("failed to get account bans: %w", err)
	}

	return bans, nil
}

func (r *AccountBanRepository) Create(ctx context.Context, input BanAccountInput) (*AccountBan, error) {
	query := `INSERT INTO account_bans (account_id, reason, banned_at, expires_at, banned_by)
	          VALUES (?, ?, UNIX_TIMESTAMP(), ?, ?)`

	_, err := r.db.ExecContext(ctx, query, input.AccountID, input.Reason, input.ExpiresAt, input.BannedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create account ban: %w", err)
	}

	bans, err := r.GetByAccountID(ctx, input.AccountID)
	if err != nil || len(bans) == 0 {
		return nil, err
	}

	return bans[0], nil
}
