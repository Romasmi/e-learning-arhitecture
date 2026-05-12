# UseCase Skill

This project uses a **Generic UseCase Pattern** with a **Type-Erased Registry** to manage business logic.

## 1. The Pattern
All business logic must implement the `UseCase[I, O]` interface:
```go
type UseCase[I, O any] interface {
    Do(ctx context.Context, input I) (O, error)
}
```

## 2. Implementing a Use Case
1. Create a struct in `internal/usecase/<domain>/`.
2. Implement the `Do` method.
3. Use domain repositories (interfaces) injected via constructor.

## 3. Registration Workflow
1. **Define ID**: Add a new `UseCaseID` to the enum in `internal/usecase/interface.go`. Update the `String()` method.
2. **Register**: Add the use case to the `Handlers` map in `internal/app/app.go` within the `registerHandlers` method:
   ```go
   a.Handlers[usecase.UseCaseMyAction] = usecase.NewHandler(myuc.NewMyUseCase(deps...))
   ```

## 4. Consumption
In any interface adapter (gRPC, CLI, etc.):
1. Get the handler: `handler := provider.GetHandler(usecase.UseCaseMyAction)`.
2. Execute: `res, err := handler.Do(ctx, inputStruct)`.
3. Cast the result: `output := res.(*domain.MyEntity)`.

## 5. Rules
- **Input/Output**: Always use structs for input and output, even if they contain only one field.
- **Errors**: Return errors as-is or wrapped with context. Use standard `fmt.Errorf("...: %w", err)`.
- **Registry**: Do not bypass the registry. Handlers should always be retrieved via `GetHandler`.
