package repository

import (
	"awesomeProject/internal/model"
	"context"
)

type Repository interface {
	Create(ctx context.Context, person *model.Person) (string, error)
	UpdateAuth(ctx context.Context, id string, refreshToken string) error
	Update(ctx context.Context, id string, person *model.Person) error
	SelectAll(ctx context.Context) ([]*model.Person, error)
	SelectByID(ctx context.Context, id string) (model.Person, error)
	Delete(ctx context.Context, id string) error
	SelectByIDAuth(ctx context.Context, id string) (model.Person, error)
}
