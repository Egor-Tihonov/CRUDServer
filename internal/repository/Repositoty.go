package repository

import (
	"awesomeProject/internal/model"
	"context"
)

type Repository interface {
	Create(ctx context.Context, person *model.Person) error
	Update(ctx context.Context, id int, person *model.Person) error
	SelectAll(ctx context.Context) ([]*model.Person, error)
	Delete(ctx context.Context, id int) error
	SelectById(ctx context.Context, id int) (model.Person, error)
}
