package models

import (
	"context"
	"fmt"

	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
)

type Account struct {
	ID            int    `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`
	Password      string `db:"password" json:"-"` // Don't expose password in JSON
	Secret        *string `db:"secret" json:"secret,omitempty"`
	Type          int    `db:"type" json:"type"`
	PremiumEndsAt int    `db:"premium_ends_at" json:"premiumEndsAt"`
	Email         string `db:"email" json:"email"`
	Creation      int    `db:"creation" json:"creation"`
}

type CreateAccountInput struct {
	Name     string
	Password string
	Email    string
}

type AccountRepository struct {
	db *database.DB
}

func NewAccountRepository(db *database.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) GetByID(ctx context.Context, id int) (*Account, error) {
	var account Account
	query := `SELECT id, name, password, secret, type, premium_ends_at, email, creation FROM accounts WHERE id = ?`

	if err := r.db.GetContext(ctx, &account, query, id); err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

func (r *AccountRepository) GetAll(ctx context.Context, limit int) ([]*Account, error) {
	var accounts []*Account
	query := `SELECT id, name, password, secret, type, premium_ends_at, email, creation FROM accounts LIMIT ?`

	if err := r.db.SelectContext(ctx, &accounts, query, limit); err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return accounts, nil
}

func (r *AccountRepository) Create(ctx context.Context, input CreateAccountInput) (*Account, error) {
	// Note: In production, hash the password with bcrypt
	query := `INSERT INTO accounts (name, password, email, creation) VALUES (?, ?, ?, UNIX_TIMESTAMP())`

	result, err := r.db.ExecContext(ctx, query, input.Name, input.Password, input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(ctx, int(id))
}
