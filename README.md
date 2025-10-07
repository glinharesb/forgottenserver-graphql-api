# Forgotten Server GraphQL API

[![CI](https://github.com/glinharesb/forgottenserver-graphql-api/workflows/CI/badge.svg)](https://github.com/glinharesb/forgottenserver-graphql-api/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/glinharesb/forgottenserver-graphql-api)](https://goreportcard.com/report/github.com/glinharesb/forgottenserver-graphql-api)
[![codecov](https://codecov.io/gh/glinharesb/forgottenserver-graphql-api/branch/main/graph/badge.svg)](https://codecov.io/gh/glinharesb/forgottenserver-graphql-api)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/glinharesb/forgottenserver-graphql-api)](https://go.dev/)

A modern, type-safe GraphQL API for [The Forgotten Server (TFS)](https://github.com/otland/forgottenserver), built with Go and gqlgen. This API provides a complete interface for managing accounts, players, guilds, houses, and marketplace operations in your OTServer.

**Compatible with:** [TFS v1.4.2](https://github.com/otland/forgottenserver/tree/v1.4.2)

## Features

- **🔐 Account Management** - Complete CRUD operations for accounts, bans, and storage
- **👥 Player System** - Player creation, statistics, deaths tracking, and storage management
- **🏰 Guild System** - Guild management with ranks, memberships, invites, and war tracking
- **🏠 House Management** - House listings, bidding system, and ownership tracking
- **💰 Market System** - Marketplace offers and transaction history
- **🏙️ Town System** - Town management and positioning
- **✅ Full Test Coverage** - Comprehensive unit tests with mocked database
- **🚀 Type-Safe** - Generated types and resolvers using gqlgen
- **📊 GraphQL Playground** - Interactive API exploration and testing

## Tech Stack

- **[Go 1.25+](https://go.dev/)** - Programming language
- **[gqlgen](https://gqlgen.com/)** - GraphQL server library and code generator
- **[MySQL](https://www.mysql.com/)** - Database (compatible with TFS schema)
- **[sqlx](https://github.com/jmoiron/sqlx)** - Enhanced SQL database driver
- **[Chi](https://github.com/go-chi/chi)** - HTTP router
- **[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)** - SQL mock driver for testing

## Prerequisites

- Go 1.25 or higher
- MySQL 5.7+ or MariaDB 10.2+
- [The Forgotten Server v1.4.2](https://github.com/otland/forgottenserver/tree/v1.4.2) database schema

## Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/glinharesb/forgottenserver-graphql-api.git
   cd forgottenserver-graphql-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up the database**

   This API requires a TFS v1.4.2 database. Set up your database using the schema from [TFS v1.4.2](https://github.com/otland/forgottenserver/tree/v1.4.2):
   ```bash
   mysql -u username -p database_name < forgottenserver/schema.sql
   ```

4. **Configure environment**

   Copy the example environment file and configure it:
   ```bash
   cp .env.example .env
   ```

   Edit `.env` with your database credentials:
   ```env
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=your_database
   SERVER_PORT=8080
   ```

5. **Run the server**
   ```bash
   make run
   ```

   Or build and run manually:
   ```bash
   go build -o server ./cmd/server
   ./server
   ```

The API will be available at `http://localhost:8080/graphql` with an interactive playground at `http://localhost:8080/`

## Project Structure

```
.
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection
│   ├── graph/           # GraphQL schema and resolvers
│   │   ├── model/       # Generated GraphQL models
│   │   └── *.graphqls   # GraphQL schema definitions
│   └── models/          # Business logic and repositories
│       ├── account.go
│       ├── player.go
│       ├── guild.go
│       ├── house.go
│       ├── market.go
│       └── ...
├── .env.example        # Example environment configuration
├── gqlgen.yml          # GraphQL code generation config
├── Makefile            # Build and development commands
└── README.md
```

## GraphQL Schema

### Queries

```graphql
type Query {
  # Accounts
  account(id: ID!): Account
  accounts(limit: Int = 10): [Account!]!

  # Players
  player(id: ID!): Player
  players(accountId: ID!): [Player!]!

  # Guilds
  guild(id: ID!): Guild
  guilds: [Guild!]!
  guildWars(guildId: ID): [GuildWar!]!

  # Houses
  house(id: ID!): House
  houses(townId: ID): [House!]!

  # Market
  marketOffers(itemType: Int): [MarketOffer!]!
  marketHistory(playerId: ID!): [MarketHistory!]!

  # Towns
  town(id: ID!): Town
  towns: [Town!]!
}
```

### Mutations

```graphql
type Mutation {
  # Accounts
  createAccount(input: CreateAccountInput!): Account!
  banAccount(input: BanAccountInput!): AccountBan!

  # Players
  createPlayer(input: CreatePlayerInput!): Player!

  # Guilds
  createGuild(input: CreateGuildInput!): Guild!
  inviteToGuild(guildId: ID!, playerId: ID!): Boolean!
  acceptGuildInvite(guildId: ID!, playerId: ID!): Boolean!

  # Houses
  bidHouse(houseId: ID!, playerId: ID!, bidAmount: Int!): House!

  # Market
  createMarketOffer(input: CreateMarketOfferInput!): MarketOffer!

  # Towns
  createTown(input: CreateTownInput!): Town!
}
```

## Example Queries

### Get Account with Players

```graphql
query GetAccount {
  account(id: "1") {
    id
    name
    email
    players {
      id
      name
      level
      vocation
    }
  }
}
```

### Create New Player

```graphql
mutation CreatePlayer {
  createPlayer(input: {
    name: "Warrior"
    accountId: 1
    sex: 1
    vocation: 4
  }) {
    id
    name
    level
    health
    healthMax
  }
}
```

### Get Guild Information

```graphql
query GetGuild {
  guild(id: "1") {
    id
    name
    owner {
      id
      name
    }
    ranks {
      id
      name
      level
    }
    members {
      player {
        name
        level
      }
      rank {
        name
      }
    }
  }
}
```

### Search Market Offers

```graphql
query GetMarketOffers {
  marketOffers(itemType: 2160) {
    id
    player {
      name
    }
    amount
    price
    created
  }
}
```

## Development

### Available Make Commands

```bash
make build          # Build the application
make run            # Run the application
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make generate       # Regenerate GraphQL code
make clean          # Clean build artifacts
make tidy           # Tidy Go modules
make docker-build   # Build Docker image
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests for specific package
go test ./internal/models/...
```

### Regenerating GraphQL Code

After modifying the GraphQL schema (`internal/graph/schema.graphqls`):

```bash
make generate
```

This will regenerate:
- Resolver interfaces
- Model types
- GraphQL execution code

## Testing

The project includes comprehensive test coverage with mocked database operations:

- **Model Tests** - Unit tests for all repository methods
- **Resolver Tests** - Integration tests for GraphQL resolvers
- **Mock Database** - Using go-sqlmock for isolated testing

Example test output:
```
✓ 27 tests in internal/graph
✓ 30 tests in internal/models
```

## Database Schema

The API is compatible with the [TFS v1.4.2](https://github.com/otland/forgottenserver/tree/v1.4.2) database schema. Key tables include:

- `accounts` - User accounts
- `players` - Player characters
- `guilds` - Guild information
- `guild_ranks` - Guild hierarchy
- `guild_membership` - Player-guild relationships
- `houses` - House data
- `market_offers` - Active market listings
- `market_history` - Completed transactions
- `towns` - Town locations
- `player_deaths` - Death history

## Configuration

Configure the application using environment variables or a `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | MySQL host | `localhost` |
| `DB_PORT` | MySQL port | `3306` |
| `DB_USER` | Database user | `root` |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | `forgottenserver` |
| `SERVER_PORT` | API server port | `8080` |

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for new features
- Follow Go best practices and conventions
- Update documentation when adding new queries/mutations
- Run `make lint` before committing
- Ensure all tests pass with `make test`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [The Forgotten Server](https://github.com/otland/forgottenserver) - The OTServer distribution this API is built for
- [gqlgen](https://gqlgen.com/) - GraphQL library for Go
- [Open Tibia](https://otland.net/) - Community and resources

## Support

- 📖 [GraphQL Playground](http://localhost:8080/) - Interactive API documentation
- 🐛 [Issue Tracker](https://github.com/glinharesb/forgottenserver-graphql-api/issues)
- 💬 [Discussions](https://github.com/glinharesb/forgottenserver-graphql-api/discussions)

---

Built with ❤️ for the Open Tibia community
