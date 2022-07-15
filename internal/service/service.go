package service

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
	"fmt"
)

type Service struct {
	rps repository.Repository
}

func NewService(NewRps repository.Repository) *Service {
	return &Service{rps: NewRps}
}

func (s *Service) CreateUser(ctx context.Context, person *model.Person) (error, string) {
	err := s.rps.Create(ctx, person)
	if err != nil {
		return fmt.Errorf("service: cannot create new user,- %v", err), ""
	}
	return nil, "Successfully create"
}

func (s *Service) UpdateUser(ctx context.Context, id int, person *model.Person) (error, string) {
	err := s.rps.Update(ctx, id, person)
	if err != nil {
		return fmt.Errorf("service: cannot to update user,- %v", err), ""
	}
	return nil, "Successfully update"
}
func (s *Service) SelectAllUsers(ctx context.Context) ([]*model.Person, error) {
	p, err := s.rps.SelectAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: cannot to select all users from database %v", err)
	}
	return p, nil
}
func (s *Service) DeleteUser(ctx context.Context, id int) (error, string) {
	err := s.rps.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service: cannot to delete user from database %v", err), ""
	}
	return nil, "Successfully delete"
}
func (s *Service) GetUserById(ctx context.Context, id int) (model.Person, error) {
	p, err := s.rps.SelectById(ctx, id)
	if err != nil {
		return p, fmt.Errorf("service: cannot to select all users from database %v", err)
	}
	return p, nil
}
