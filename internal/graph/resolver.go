package graph

import (
	"github.com/glinharesb/forgottenserver-graphql-api/internal/database"
	"github.com/glinharesb/forgottenserver-graphql-api/internal/models"
)

// Resolver is the root GraphQL resolver
type Resolver struct {
	DB                       *database.DB
	AccountRepository        *models.AccountRepository
	AccountBanRepository     *models.AccountBanRepository
	AccountStorageRepository *models.AccountStorageRepository
	PlayerRepository         *models.PlayerRepository
	PlayerDeathRepository    *models.PlayerDeathRepository
	PlayerStorageRepository  *models.PlayerStorageRepository
	TownRepository           *models.TownRepository
	GuildRepository          *models.GuildRepository
	HouseRepository          *models.HouseRepository
	MarketRepository         *models.MarketRepository
}

func NewResolver(db *database.DB) *Resolver {
	return &Resolver{
		DB:                       db,
		AccountRepository:        models.NewAccountRepository(db),
		AccountBanRepository:     models.NewAccountBanRepository(db),
		AccountStorageRepository: models.NewAccountStorageRepository(db),
		PlayerRepository:         models.NewPlayerRepository(db),
		PlayerDeathRepository:    models.NewPlayerDeathRepository(db),
		PlayerStorageRepository:  models.NewPlayerStorageRepository(db),
		TownRepository:           models.NewTownRepository(db),
		GuildRepository:          models.NewGuildRepository(db),
		HouseRepository:          models.NewHouseRepository(db),
		MarketRepository:         models.NewMarketRepository(db),
	}
}
