package models

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuildRepository_GetByID(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewGuildRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "ownerid", "creationdata", "motd"}).
		AddRow(1, "Test Guild", 1, 1234567890, "Welcome!")

	mock.ExpectQuery("SELECT id, name, ownerid, creationdata, motd FROM guilds WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	guild, err := repo.GetByID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 1, guild.ID)
	assert.Equal(t, "Test Guild", guild.Name)
	assert.Equal(t, 1, guild.OwnerID)
	assert.Equal(t, "Welcome!", guild.MOTD)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGuildRepository_GetAll(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewGuildRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "ownerid", "creationdata", "motd"}).
		AddRow(1, "Guild 1", 1, 1234567890, "").
		AddRow(2, "Guild 2", 2, 1234567891, "")

	mock.ExpectQuery("SELECT id, name, ownerid, creationdata, motd FROM guilds").
		WillReturnRows(rows)

	guilds, err := repo.GetAll(context.Background())

	require.NoError(t, err)
	assert.Len(t, guilds, 2)
	assert.Equal(t, "Guild 1", guilds[0].Name)
	assert.Equal(t, "Guild 2", guilds[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGuildRepository_Create(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewGuildRepository(db)

	input := CreateGuildInput{
		Name:    "New Guild",
		OwnerID: 1,
	}

	mock.ExpectExec("INSERT INTO guilds").
		WithArgs(input.Name, input.OwnerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id", "name", "ownerid", "creationdata", "motd"}).
		AddRow(1, input.Name, input.OwnerID, 1234567890, "")

	mock.ExpectQuery("SELECT id, name, ownerid, creationdata, motd FROM guilds WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	guild, err := repo.Create(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, "New Guild", guild.Name)
	assert.Equal(t, 1, guild.OwnerID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGuildRepository_GetRanks(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewGuildRepository(db)

	rows := sqlmock.NewRows([]string{"id", "guild_id", "name", "level"}).
		AddRow(1, 1, "Leader", 3).
		AddRow(2, 1, "Vice-Leader", 2).
		AddRow(3, 1, "Member", 1)

	mock.ExpectQuery("SELECT id, guild_id, name, level FROM guild_ranks WHERE guild_id = \\? ORDER BY level DESC").
		WithArgs(1).
		WillReturnRows(rows)

	ranks, err := repo.GetRanks(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, ranks, 3)
	assert.Equal(t, "Leader", ranks[0].Name)
	assert.Equal(t, 3, ranks[0].Level)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGuildRepository_GetMembers(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewGuildRepository(db)

	rows := sqlmock.NewRows([]string{"player_id", "guild_id", "rank_id", "nick"}).
		AddRow(1, 1, 1, "Leader Nick").
		AddRow(2, 1, 2, "Member Nick")

	mock.ExpectQuery("SELECT player_id, guild_id, rank_id, nick FROM guild_membership WHERE guild_id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	members, err := repo.GetMembers(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, members, 2)
	assert.Equal(t, "Leader Nick", members[0].Nick)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGuildRepository_GetMembershipByPlayerID(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewGuildRepository(db)

	rows := sqlmock.NewRows([]string{"player_id", "guild_id", "rank_id", "nick"}).
		AddRow(1, 1, 1, "Leader Nick")

	mock.ExpectQuery("SELECT player_id, guild_id, rank_id, nick FROM guild_membership WHERE player_id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	membership, err := repo.GetMembershipByPlayerID(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, 1, membership.PlayerID)
	assert.Equal(t, 1, membership.GuildID)
	assert.Equal(t, "Leader Nick", membership.Nick)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGuildRepository_InvitePlayer(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewGuildRepository(db)

	mock.ExpectExec("INSERT INTO guild_invites").
		WithArgs(2, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.InvitePlayer(context.Background(), 1, 2)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGuildRepository_AcceptInvite(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewGuildRepository(db)

	// Mock getting rank
	rankRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT id FROM guild_ranks WHERE guild_id = \\? AND level = 1 LIMIT 1").
		WithArgs(1).
		WillReturnRows(rankRows)

	// Mock adding to guild_membership
	mock.ExpectExec("INSERT INTO guild_membership").
		WithArgs(2, 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Mock removing invite
	mock.ExpectExec("DELETE FROM guild_invites WHERE player_id = \\? AND guild_id = \\?").
		WithArgs(2, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.AcceptInvite(context.Background(), 1, 2)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGuildRepository_GetWars(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewGuildRepository(db)

	t.Run("GetAllWars", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "guild1", "guild2", "name1", "name2", "status", "started", "ended"}).
			AddRow(1, 1, 2, "Guild1", "Guild2", 1, 1234567890, 0)

		mock.ExpectQuery("SELECT id, guild1, guild2, name1, name2, status, started, ended FROM guild_wars").
			WillReturnRows(rows)

		wars, err := repo.GetWars(context.Background(), nil)

		require.NoError(t, err)
		assert.Len(t, wars, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetWarsForGuild", func(t *testing.T) {
		guildID := 1
		rows := sqlmock.NewRows([]string{"id", "guild1", "guild2", "name1", "name2", "status", "started", "ended"}).
			AddRow(1, 1, 2, "Guild1", "Guild2", 1, 1234567890, 0)

		mock.ExpectQuery("SELECT id, guild1, guild2, name1, name2, status, started, ended FROM guild_wars WHERE guild1 = \\? OR guild2 = \\?").
			WithArgs(guildID, guildID).
			WillReturnRows(rows)

		wars, err := repo.GetWars(context.Background(), &guildID)

		require.NoError(t, err)
		assert.Len(t, wars, 1)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGuildRepository_GetWarKills(t *testing.T) {
	db, mock, err := NewMockDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewGuildRepository(db)

	rows := sqlmock.NewRows([]string{"id", "killer", "target", "killerguild", "targetguild", "warid", "time"}).
		AddRow(1, "Player1", "Player2", 1, 2, 1, 1234567890)

	mock.ExpectQuery("SELECT id, killer, target, killerguild, targetguild, warid, time FROM guildwar_kills WHERE warid = \\? ORDER BY time DESC").
		WithArgs(1).
		WillReturnRows(rows)

	kills, err := repo.GetWarKills(context.Background(), 1)

	require.NoError(t, err)
	assert.Len(t, kills, 1)
	assert.Equal(t, "Player1", kills[0].Killer)
	assert.Equal(t, "Player2", kills[0].Target)
	assert.NoError(t, mock.ExpectationsWereMet())
}
