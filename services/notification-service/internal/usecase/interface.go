package usecase

import (
	"context"
	"fmt"
)

// UseCaseID identifies a use case in the registry.
type UseCaseID int

const (
	UseCaseUnknown UseCaseID = iota
	UseCaseListMessages
)

func (id UseCaseID) String() string {
	switch id {
	case UseCaseListMessages:
		return "ListMessages"
	default:
		return "Unknown"
	}
}

// UseCase defines the generic interface for a business logic operation.
type UseCase[I, O any] interface {
	Do(ctx context.Context, input I) (O, error)
}

// Handler is a type-erased interface for use cases, useful for registries.
type Handler interface {
	Do(ctx context.Context, input any) (any, error)
}

// handlerAdapter wraps a typed UseCase into a type-erased Handler.
type handlerAdapter[I, O any] struct {
	uc UseCase[I, O]
}

func (h *handlerAdapter[I, O]) Do(ctx context.Context, input any) (any, error) {
	typed, ok := input.(I)
	if !ok {
		return nil, fmt.Errorf("invalid input type: expected %T, got %T", *new(I), input)
	}
	return h.uc.Do(ctx, typed)
}

// NewHandler creates a new Handler from a UseCase.
func NewHandler[I, O any](uc UseCase[I, O]) Handler {
	return &handlerAdapter[I, O]{uc: uc}
}
