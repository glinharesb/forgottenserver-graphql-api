package models

import (
	"context"
	"fmt"

	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
)

type Guild struct {
	ID           int    `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	OwnerID      int    `db:"ownerid" json:"ownerId"`
	CreationData int    `db:"creationdata" json:"creationData"`
	MOTD         string `db:"motd" json:"motd"`
}

type GuildRank struct {
	ID      int    `db:"id" json:"id"`
	GuildID int    `db:"guild_id" json:"guildId"`
	Name    string `db:"name" json:"name"`
	Level   int    `db:"level" json:"level"`
}

type GuildMembership struct {
	PlayerID int    `db:"player_id" json:"playerId"`
	GuildID  int    `db:"guild_id" json:"guildId"`
	RankID   int    `db:"rank_id" json:"rankId"`
	Nick     string `db:"nick" json:"nick"`
}

type GuildInvite struct {
	PlayerID int `db:"player_id" json:"playerId"`
	GuildID  int `db:"guild_id" json:"guildId"`
}

type GuildWar struct {
	ID      int    `db:"id" json:"id"`
	Guild1  int    `db:"guild1" json:"guild1"`
	Guild2  int    `db:"guild2" json:"guild2"`
	Name1   string `db:"name1" json:"name1"`
	Name2   string `db:"name2" json:"name2"`
	Status  int    `db:"status" json:"status"`
	Started int64  `db:"started" json:"started"`
	Ended   int64  `db:"ended" json:"ended"`
}

type GuildWarKill struct {
	ID          int    `db:"id" json:"id"`
	Killer      string `db:"killer" json:"killer"`
	Target      string `db:"target" json:"target"`
	KillerGuild int    `db:"killerguild" json:"killerGuild"`
	TargetGuild int    `db:"targetguild" json:"targetGuild"`
	WarID       int    `db:"warid" json:"warId"`
	Time        int64  `db:"time" json:"time"`
}

type CreateGuildInput struct {
	Name    string
	OwnerID int
}

type GuildRepository struct {
	db *database.DB
}

func NewGuildRepository(db *database.DB) *GuildRepository {
	return &GuildRepository{db: db}
}

func (r *GuildRepository) GetByID(ctx context.Context, id int) (*Guild, error) {
	var guild Guild
	query := `SELECT id, name, ownerid, creationdata, motd FROM guilds WHERE id = ?`

	if err := r.db.GetContext(ctx, &guild, query, id); err != nil {
		return nil, fmt.Errorf("failed to get guild: %w", err)
	}

	return &guild, nil
}

func (r *GuildRepository) GetAll(ctx context.Context) ([]*Guild, error) {
	var guilds []*Guild
	query := `SELECT id, name, ownerid, creationdata, motd FROM guilds`

	if err := r.db.SelectContext(ctx, &guilds, query); err != nil {
		return nil, fmt.Errorf("failed to get guilds: %w", err)
	}

	return guilds, nil
}

func (r *GuildRepository) Create(ctx context.Context, input CreateGuildInput) (*Guild, error) {
	query := `INSERT INTO guilds (name, ownerid, creationdata) VALUES (?, ?, UNIX_TIMESTAMP())`

	result, err := r.db.ExecContext(ctx, query, input.Name, input.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to create guild: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(ctx, int(id))
}

func (r *GuildRepository) GetRanks(ctx context.Context, guildID int) ([]*GuildRank, error) {
	var ranks []*GuildRank
	query := `SELECT id, guild_id, name, level FROM guild_ranks WHERE guild_id = ? ORDER BY level DESC`

	if err := r.db.SelectContext(ctx, &ranks, query, guildID); err != nil {
		return nil, fmt.Errorf("failed to get guild ranks: %w", err)
	}

	return ranks, nil
}

func (r *GuildRepository) GetMembers(ctx context.Context, guildID int) ([]*GuildMembership, error) {
	var members []*GuildMembership
	query := `SELECT player_id, guild_id, rank_id, nick FROM guild_membership WHERE guild_id = ?`

	if err := r.db.SelectContext(ctx, &members, query, guildID); err != nil {
		return nil, fmt.Errorf("failed to get guild members: %w", err)
	}

	return members, nil
}

func (r *GuildRepository) GetMembershipByPlayerID(ctx context.Context, playerID int) (*GuildMembership, error) {
	var membership GuildMembership
	query := `SELECT player_id, guild_id, rank_id, nick FROM guild_membership WHERE player_id = ?`

	if err := r.db.GetContext(ctx, &membership, query, playerID); err != nil {
		return nil, fmt.Errorf("failed to get guild membership: %w", err)
	}

	return &membership, nil
}

func (r *GuildRepository) InvitePlayer(ctx context.Context, guildID, playerID int) error {
	query := `INSERT INTO guild_invites (player_id, guild_id) VALUES (?, ?)`

	_, err := r.db.ExecContext(ctx, query, playerID, guildID)
	if err != nil {
		return fmt.Errorf("failed to invite player: %w", err)
	}

	return nil
}

func (r *GuildRepository) AcceptInvite(ctx context.Context, guildID, playerID int) error {
	// Get lowest rank (level 1)
	var rankID int
	query := `SELECT id FROM guild_ranks WHERE guild_id = ? AND level = 1 LIMIT 1`
	if err := r.db.GetContext(ctx, &rankID, query, guildID); err != nil {
		return fmt.Errorf("failed to get rank: %w", err)
	}

	// Add to guild_membership
	query = `INSERT INTO guild_membership (player_id, guild_id, rank_id) VALUES (?, ?, ?)`
	if _, err := r.db.ExecContext(ctx, query, playerID, guildID, rankID); err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	// Remove invite
	query = `DELETE FROM guild_invites WHERE player_id = ? AND guild_id = ?`
	if _, err := r.db.ExecContext(ctx, query, playerID, guildID); err != nil {
		return fmt.Errorf("failed to remove invite: %w", err)
	}

	return nil
}

func (r *GuildRepository) GetWars(ctx context.Context, guildID *int) ([]*GuildWar, error) {
	var wars []*GuildWar
	var query string

	if guildID != nil {
		query = `SELECT id, guild1, guild2, name1, name2, status, started, ended
		         FROM guild_wars WHERE guild1 = ? OR guild2 = ?`
		if err := r.db.SelectContext(ctx, &wars, query, *guildID, *guildID); err != nil {
			return nil, fmt.Errorf("failed to get guild wars: %w", err)
		}
	} else {
		query = `SELECT id, guild1, guild2, name1, name2, status, started, ended FROM guild_wars`
		if err := r.db.SelectContext(ctx, &wars, query); err != nil {
			return nil, fmt.Errorf("failed to get guild wars: %w", err)
		}
	}

	return wars, nil
}

func (r *GuildRepository) GetWarKills(ctx context.Context, warID int) ([]*GuildWarKill, error) {
	var kills []*GuildWarKill
	query := `SELECT id, killer, target, killerguild, targetguild, warid, time
	          FROM guildwar_kills WHERE warid = ? ORDER BY time DESC`

	if err := r.db.SelectContext(ctx, &kills, query, warID); err != nil {
		return nil, fmt.Errorf("failed to get war kills: %w", err)
	}

	return kills, nil
}
