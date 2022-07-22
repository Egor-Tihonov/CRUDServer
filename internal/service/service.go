package service

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
)

var JwtKey = []byte("super-key") //key fo generation and check tokens

type Service struct { //service new
	rps repository.Repository
}

func NewService(NewRps repository.Repository) *Service { //create
	return &Service{rps: NewRps}
}

func (s *Service) UpdateUser(ctx context.Context, id string, person *model.Person) error { //update user
	return s.rps.Update(ctx, id, person)
}
func (s *Service) SelectAllUsers(ctx context.Context) ([]*model.Person, error) { //get all users from DB without passwords and tokens
	return s.rps.SelectAll(ctx)
}
func (s *Service) DeleteUser(ctx context.Context, id string) error { //delete user from DB
	return s.rps.Delete(ctx, id)
}
func (s *Service) GetUserById(ctx context.Context, id string) (model.Person, error) { //get one user by id
	return s.rps.SelectById(ctx, id)
}
