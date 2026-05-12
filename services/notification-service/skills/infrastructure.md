# Infrastructure Skill

Infrastructure layer handles technical implementation details for data storage and external communication.

## 1. Database (PostgreSQL)
- **Driver**: `pgx/v5`.
- **Repository**: Implement interfaces from the `domain` layer in `internal/infrastructure/db/postgres/`.
- **Migrations**: Stored in `migrations/`, managed via `migrate` tool (see `Makefile`).
- **Connection**: Managed via `pgxpool.Pool` initialized in `internal/app/app.go`.

## 2. Kafka
- **Driver**: `segmentio/kafka-go`.
- **Producer**: Generic or domain-specific writers in `internal/infrastructure/kafka/`.
- **Consumer**: Handlers in `internal/interface/kafka/`.

## 3. Configuration
- **Library**: Viper.
- **Files**: `config.yaml` and environment variables.
- **Implementation**: `internal/config/config.go`.
- *Rule*: New configuration fields must be added to the `Config` struct and properly tagged with `mapstructure`.
