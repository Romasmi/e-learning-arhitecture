# Architecture Skill

This project follows **Clean Architecture** and **DDD (Domain-Driven Design)** principles.

## 1. Principles
- **Dependency Rule**: Dependencies point inwards. Outer layers (Interface, Infrastructure) depend on inner layers (UseCase, Domain). Domain layer depends on nothing.
- **Separation of Concerns**: Each layer has a specific responsibility.
- **Pure Domain**: The domain layer contains only business logic and entities, no technical details.

## 2. Layers

### Domain (`internal/domain/`)
- **Entities**: Core business objects (e.g., `user.User`).
- **Interfaces**: Repository or Service interfaces that the domain needs.
- *Rule*: No imports from other layers.

### UseCase (`internal/usecase/`)
- **Application Logic**: Orchestrates the flow of data to and from entities.
- **Implementation**: Every operation must implement `UseCase[I, O]`.
- *Rule*: Depends only on Domain.

### Infrastructure (`internal/infrastructure/`)
- **Technical Implementations**: Database (PostgreSQL), Messaging (Kafka), External APIs.
- *Rule*: Implements interfaces defined in Domain or UseCase.

### Interface (`internal/interface/`)
- **Adapters**: gRPC handlers, HTTP/REST Gateway, CLI commands, Kafka Consumers.
- *Rule*: Converts external requests into UseCase inputs and calls UseCase handlers.

## 3. Entry Points (`cmd/`)
- `api/`: Runs gRPC and HTTP Gateway.
- `worker/`: Runs Kafka consumers and background tasks.
- `cli/`: Management commands via Cobra.
