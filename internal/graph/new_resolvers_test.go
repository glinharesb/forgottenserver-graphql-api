package graph

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/glinharesb/forgottenserver-graphql-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Town Query Tests

func TestQueryResolver_Town(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "posx", "posy", "posz"}).
		AddRow(1, "Thais", 100, 200, 7)

	mock.ExpectQuery("SELECT id, name, posx, posy, posz FROM towns WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	town, err := resolver.Query().Town(context.Background(), "1")

	require.NoError(t, err)
	assert.Equal(t, 1, town.ID)
	assert.Equal(t, "Thais", town.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryResolver_Towns(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "posx", "posy", "posz"}).
		AddRow(1, "Thais", 100, 200, 7).
		AddRow(2, "Carlin", 150, 250, 7)

	mock.ExpectQuery("SELECT id, name, posx, posy, posz FROM towns").
		WillReturnRows(rows)

	towns, err := resolver.Query().Towns(context.Background())

	require.NoError(t, err)
	assert.Len(t, towns, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Guild Query Tests

func TestQueryResolver_Guild(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "ownerid", "creationdata", "motd"}).
		AddRow(1, "Test Guild", 1, 1234567890, "Welcome!")

	mock.ExpectQuery("SELECT id, name, ownerid, creationdata, motd FROM guilds WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	guild, err := resolver.Query().Guild(context.Background(), "1")

	require.NoError(t, err)
	assert.Equal(t, 1, guild.ID)
	assert.Equal(t, "Test Guild", guild.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryResolver_Guilds(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "ownerid", "creationdata", "motd"}).
		AddRow(1, "Guild 1", 1, 1234567890, "").
		AddRow(2, "Guild 2", 2, 1234567891, "")

	mock.ExpectQuery("SELECT id, name, ownerid, creationdata, motd FROM guilds").
		WillReturnRows(rows)

	guilds, err := resolver.Query().Guilds(context.Background())

	require.NoError(t, err)
	assert.Len(t, guilds, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// House Query Tests

func TestQueryResolver_House(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "owner", "paid", "warnings", "name", "rent", "town_id", "bid", "bid_end",
		"last_bid", "highest_bidder", "size", "beds",
	}).AddRow(1, 0, 0, 0, "Test House", 1000, 1, 0, 0, 0, 0, 100, 2)

	mock.ExpectQuery("SELECT (.+) FROM houses WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	house, err := resolver.Query().House(context.Background(), "1")

	require.NoError(t, err)
	assert.Equal(t, 1, house.ID)
	assert.Equal(t, "Test House", house.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Market Query Tests

func TestQueryResolver_MarketOffers(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "player_id", "sale", "itemtype", "amount", "created", "anonymous", "price",
	}).AddRow(1, 1, true, 2160, 10, 1234567890, false, 1000)

	itemType := 2160
	mock.ExpectQuery("SELECT (.+) FROM market_offers WHERE itemtype = \\? ORDER BY created DESC").
		WithArgs(itemType).
		WillReturnRows(rows)

	offers, err := resolver.Query().MarketOffers(context.Background(), &itemType)

	require.NoError(t, err)
	assert.Len(t, offers, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestQueryResolver_MarketHistory(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{
		"id", "player_id", "sale", "itemtype", "amount", "price", "expires_at", "inserted", "state",
	}).AddRow(1, 1, true, 2160, 10, 1000, 1234567890, 1234567800, 1)

	mock.ExpectQuery("SELECT (.+) FROM market_history WHERE player_id = \\? ORDER BY inserted DESC").
		WithArgs(1).
		WillReturnRows(rows)

	history, err := resolver.Query().MarketHistory(context.Background(), "1")

	require.NoError(t, err)
	assert.Len(t, history, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Mutation Tests

func TestMutationResolver_CreateTown(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	input := models.CreateTownInput{
		Name: "NewTown",
		PosX: 300,
		PosY: 400,
		PosZ: 7,
	}

	mock.ExpectExec("INSERT INTO towns").
		WithArgs(input.Name, input.PosX, input.PosY, input.PosZ).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id", "name", "posx", "posy", "posz"}).
		AddRow(1, input.Name, input.PosX, input.PosY, input.PosZ)

	mock.ExpectQuery("SELECT id, name, posx, posy, posz FROM towns WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	town, err := resolver.Mutation().CreateTown(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, "NewTown", town.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMutationResolver_CreateGuild(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	input := models.CreateGuildInput{
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

	guild, err := resolver.Mutation().CreateGuild(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, "New Guild", guild.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMutationResolver_BanAccount(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	input := models.BanAccountInput{
		AccountID: 1,
		Reason:    "Botting",
		ExpiresAt: 1234567900,
		BannedBy:  1,
	}

	mock.ExpectExec("INSERT INTO account_bans").
		WithArgs(input.AccountID, input.Reason, input.ExpiresAt, input.BannedBy).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rows := sqlmock.NewRows([]string{"account_id", "reason", "banned_at", "expires_at", "banned_by"}).
		AddRow(input.AccountID, input.Reason, 1234567890, input.ExpiresAt, input.BannedBy)

	mock.ExpectQuery("SELECT account_id, reason, banned_at, expires_at, banned_by FROM account_bans WHERE account_id = ?").
		WithArgs(input.AccountID).
		WillReturnRows(rows)

	ban, err := resolver.Mutation().BanAccount(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, "Botting", ban.Reason)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMutationResolver_CreateMarketOffer(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	input := models.CreateMarketOfferInput{
		PlayerID:  1,
		Sale:      true,
		ItemType:  2160,
		Amount:    10,
		Price:     1000,
		Anonymous: false,
	}

	mock.ExpectExec("INSERT INTO market_offers").
		WithArgs(input.PlayerID, input.Sale, input.ItemType, input.Amount, input.Anonymous, input.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{
		"id", "player_id", "sale", "itemtype", "amount", "created", "anonymous", "price",
	}).AddRow(1, input.PlayerID, input.Sale, input.ItemType, input.Amount, 1234567890, input.Anonymous, input.Price)

	mock.ExpectQuery("SELECT (.+) FROM market_offers WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	offer, err := resolver.Mutation().CreateMarketOffer(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, 2160, offer.ItemType)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Field Resolver Tests

func TestPlayerResolver_Town(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	player := &models.Player{
		ID:     1,
		TownID: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "name", "posx", "posy", "posz"}).
		AddRow(1, "Thais", 100, 200, 7)

	mock.ExpectQuery("SELECT id, name, posx, posy, posz FROM towns WHERE id = ?").
		WithArgs(player.TownID).
		WillReturnRows(rows)

	town, err := resolver.Player().Town(context.Background(), player)

	require.NoError(t, err)
	assert.Equal(t, "Thais", town.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPlayerResolver_Deaths(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	player := &models.Player{
		ID: 1,
	}

	rows := sqlmock.NewRows([]string{
		"player_id", "time", "level", "killed_by", "is_player", "mostdamage_by", "mostdamage_is_player",
	}).AddRow(1, 1234567890, 50, "Dragon", false, "Dragon", false)

	mock.ExpectQuery("SELECT (.+) FROM player_deaths WHERE player_id = \\? ORDER BY time DESC").
		WithArgs(player.ID).
		WillReturnRows(rows)

	deaths, err := resolver.Player().Deaths(context.Background(), player)

	require.NoError(t, err)
	assert.Len(t, deaths, 1)
	assert.Equal(t, "Dragon", deaths[0].KilledBy)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountResolver_Bans(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	account := &models.Account{
		ID: 1,
	}

	rows := sqlmock.NewRows([]string{"account_id", "reason", "banned_at", "expires_at", "banned_by"}).
		AddRow(1, "Botting", 1234567890, 1234567900, 1)

	mock.ExpectQuery("SELECT account_id, reason, banned_at, expires_at, banned_by FROM account_bans WHERE account_id = ?").
		WithArgs(account.ID).
		WillReturnRows(rows)

	bans, err := resolver.Account().Bans(context.Background(), account)

	require.NoError(t, err)
	assert.Len(t, bans, 1)
	assert.Equal(t, "Botting", bans[0].Reason)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGuildResolver_Owner(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	guild := &models.Guild{
		ID:      1,
		OwnerID: 1,
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "group_id", "account_id", "level", "vocation", "health", "healthmax",
		"experience", "lookbody", "lookfeet", "lookhead", "looklegs", "looktype", "lookaddons",
		"maglevel", "mana", "manamax", "soul", "town_id", "posx", "posy", "posz", "cap", "sex",
		"lastlogin", "balance",
	}).AddRow(1, "Owner", 1, 1, 100, 4, 1000, 1000, 500000, 0, 0, 0, 0, 136, 0, 50, 500, 500, 0, 1, 0, 0, 0, 500, 1, 0, 10000)

	mock.ExpectQuery("SELECT (.+) FROM players WHERE id = ?").
		WithArgs(guild.OwnerID).
		WillReturnRows(rows)

	owner, err := resolver.Guild().Owner(context.Background(), guild)

	require.NoError(t, err)
	assert.Equal(t, "Owner", owner.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGuildResolver_Ranks(t *testing.T) {
	resolver, mock, cleanup := setupTestResolver(t)
	defer cleanup()

	guild := &models.Guild{
		ID: 1,
	}

	rows := sqlmock.NewRows([]string{"id", "guild_id", "name", "level"}).
		AddRow(1, 1, "Leader", 3).
		AddRow(2, 1, "Member", 1)

	mock.ExpectQuery("SELECT id, guild_id, name, level FROM guild_ranks WHERE guild_id = \\? ORDER BY level DESC").
		WithArgs(guild.ID).
		WillReturnRows(rows)

	ranks, err := resolver.Guild().Ranks(context.Background(), guild)

	require.NoError(t, err)
	assert.Len(t, ranks, 2)
	assert.Equal(t, "Leader", ranks[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}
