package domain

import "context"

type Repository[T any] interface {
	Save(ctx context.Context, entity *T) error
	GetById(ctx context.Context, id string) (*T, error)
	GetAll(ctx context.Context) ([]*T, error)
}
