package repository

import (
	"awesomeProject/internal/model"
	"context"
)

type Repository interface {
	Create(ctx context.Context, person *model.Person) (error, string)
	Update(ctx context.Context, id string, person *model.Person) error
	SelectAll(ctx context.Context) ([]*model.Person, error)
	Delete(ctx context.Context, id string) error
	SelectById(ctx context.Context, id string) (model.Person, error)
}
