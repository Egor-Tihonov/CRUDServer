package service

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
)

type Service struct {
	rps repository.Repository
}

func NewService(NewRps repository.Repository) *Service {
	return &Service{rps: NewRps}
}

func (s *Service) CreateUser(ctx context.Context, person *model.Person) (error, string) {
	return s.rps.Create(ctx, person)
}

func (s *Service) UpdateUser(ctx context.Context, id string, person *model.Person) error {
	return s.rps.Update(ctx, id, person)
}
func (s *Service) SelectAllUsers(ctx context.Context) ([]*model.Person, error) {
	return s.rps.SelectAll(ctx)
}
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	return s.rps.Delete(ctx, id)
}
func (s *Service) GetUserById(ctx context.Context, id string) (model.Person, error) {
	return s.rps.SelectById(ctx, id)
}
