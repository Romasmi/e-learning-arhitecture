# Domain Skill

The Domain layer is the core of the application. It contains the business entities and the logic that governs them.

## 1. Entities
- Entities represent business objects with identity.
- Example: `internal/domain/user/entity.go`.
- *Rule*: Use basic Go types or standard libraries. Avoid external dependencies in entity definitions.

## 2. Repository Interfaces
- Define how entities are stored and retrieved.
- Interfaces are defined in the domain package but implemented in the `infrastructure` layer.
- *Rule*: Repository methods should accept and return domain entities.

## 3. Pure Domain Logic
- Complex business rules that don't belong to a single entity should be implemented as Domain Services in the domain package.

## 4. Constraint: No Inward Dependencies
- The Domain layer must **never** import from `usecase`, `interface`, or `infrastructure`.
