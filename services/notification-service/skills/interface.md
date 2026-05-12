# Interface Skill

Interface adapters convert external signals (gRPC, HTTP, CLI, Kafka) into calls to the UseCase layer.

## 1. gRPC
- **Files**: `internal/interface/grpc/`.
- **Logic**: Use `UseCaseProvider` to get handlers.
- **Registration**: Add service registration to `internal/interface/grpc/server.go`.
- **Mapping**: Convert Protobuf messages to UseCase input structs and back.

## 2. CLI
- **Framework**: Cobra.
- **Files**: `internal/interface/cli/`.
- **Dependency**: Uses the global `deps` (set via `SetApp`) which implements `AppDependencies`.
- **Logic**: Commands retrieve handlers from `deps` and execute them.

## 3. Kafka
- **Producer**: Implement domain events in `internal/infrastructure/kafka/`.
- **Consumer**: Implement in `internal/interface/kafka/`.
- **Lifecycle**: Consumers must implement `Start(ctx)` and `Close()`.
- **Registration**: Start consumers in `internal/app/worker.go`.

## 4. HTTP Gateway
- **Config**: Defined in `.proto` files.
- **Server**: `internal/interface/http/server.go` sets up the `runtime.NewServeMux` and handles documentation routes (`/swagger`, `/proto`).
