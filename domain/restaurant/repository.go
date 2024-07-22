package restaurant

import (
	"context"
	"github.com/google/uuid"
)

type Repository[T comparable] interface {
	Create(ctx context.Context, o T) error
	Get(ctx context.Context, id uuid.UUID) (T, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
