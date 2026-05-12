# Go Microservice Template

A modern Go microservice template following DDD and Clean Architecture principles.

## Features

- **gRPC & gRPC Gateway**: Primary communication via gRPC with REST support.
- **Kafka**: Integrated event-driven communication (Producer/Consumer).
- **DDD & Clean Architecture**: Organized into Domain, Usecase, Infrastructure, and Interface layers.
- **Multiple Entrypoints**:
  - `cmd/api`: Runs gRPC and HTTP Gateway servers.
  - `cmd/worker`: Runs background Kafka consumers.
  - `cmd/cli`: Command-line interface using Cobra.
- **Documentation**:
  - Swagger UI at `/swagger`
  - Proto schema at `/proto`
- **Environment Configuration**: Managed via Viper (supports `.yaml` and env vars).
- **Local Development**: Docker Compose for DB and Kafka.
- **Developer & AI Skills**: Specialized documentation for AI agents (Junie, Cursor) and developers in the `skills/` and `.junie/` directories.

## Project Structure

```text
api/                    # Protocol Buffer definitions
cmd/                    # Application entry points
  api/                  # gRPC & Gateway server
  worker/               # Kafka consumer
  cli/                  # Command line interface
configs/                # Configuration files
deployments/            # Docker Compose and K8s manifests
internal/
  api/                  # Generated gRPC & Gateway code
  app/                  # Application bootstrap
  config/               # Configuration loading
  domain/               # Domain entities and interfaces (DDD)
  usecase/              # Business logic (Application Services)
  infrastructure/       # External implementations (DB, Kafka)
  interface/            # Adapters (gRPC, HTTP, Kafka Handlers)
migrations/             # Database migrations
skills/                 # AI & Developer skills (patterns and rules)
.junie/                 # Junie-specific AI guidelines
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- `buf` (for code generation)

### Setup

1. **Clone the repository**:
   ```bash
   git clone <repo-url>
   cd go-api-template
   ```

2. **Start local infrastructure**:
   ```bash
   make compose-up
   ```

3. **Generate code**:
   ```bash
   make generate
   ```

4. **Run migrations**:
   ```bash
   make migrate-up
   ```

### Running the services

- **API Server**: `make run-api`
- **Worker**: `make run-worker`
- **CLI**: `make run-cli user create "John Doe" "john@example.com"`

## How to Add a New Feature

1. **Define Proto**: Add a new `.proto` file in `api/` or update existing ones.
2. **Generate Code**: Run `make generate`.
3. **Domain Layer**: Define your entity in `internal/domain/`.
4. **Infrastructure Layer**: Implement the repository interface in `internal/infrastructure/db/postgres/`.
5. **Usecase Layer**: 
   - Implement the business logic in a new struct in `internal/usecase/your_domain/`.
   - The struct should implement the `usecase.UseCase[I, O]` interface with a `Do` method.
6. **Register Usecase**: 
   - Add a new entry to the `UseCaseID` enum in `internal/usecase/interface.go`.
   - Register your usecase in `internal/app/app.go` within the `registerHandlers` method by adding it to the `Handlers` map using `usecase.NewHandler`.
7. **Interface Layer**: 
   - Use the registered handler in gRPC handler (`internal/interface/grpc/`) by calling `app.GetHandler(usecase.YourUseCaseID)`.
   - (Optional) Register it in gRPC gateway in `internal/interface/http/server.go`.
   - (Optional) Add Kafka consumer in `internal/interface/kafka/` and register in `internal/app/worker.go`.
   - (Optional) Add CLI command in `internal/interface/cli/` and register in `internal/app/cli.go`.

## API Documentation

- **Swagger UI**: Rendered UI is available at `http://localhost:8080/swagger/`. It uses the generated Swagger JSON files.
- **Proto Files**: Served at `http://localhost:8080/proto/`.

## Event-Driven Flow Example

1. `cmd/api` receives a `CreateUser` request.
2. `usecase` saves the user to PostgreSQL via `UserRepository`.
3. `usecase` publishes a `UserCreated` event via `EventProducer`.
4. `cmd/worker` (the consumer app) receives the event and logs it (or performs other actions like sending an email).

## CLI Examples

The CLI provides management commands. 

```bash
# Create a user
go run cmd/cli/main.go user create "Alice" "alice@example.com"

# Reset password (example command)
go run cmd/cli/main.go user reset-password "alice@example.com"
```

## DDD and Clean Architecture

- **Domain**: Contains entities and repository interfaces. No external dependencies.
- **Usecase**: Implements business logic using domain entities and interfaces.
- **Infrastructure**: Implements repository interfaces using specific technologies (PostgreSQL, Kafka).
- **Interface**: Adapts external requests (gRPC, HTTP, Kafka) to usecases.
