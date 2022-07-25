package service

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
	"fmt"
)

var JwtKey = []byte("super-key") //key fo generation and check tokens

type Service struct { //service new
	rps repository.Repository
}

func NewService(NewRps repository.Repository) *Service { //create
	return &Service{rps: NewRps}
}

func (s *Service) UpdateUser(ctx context.Context, id string, person *model.Person) error { //update user
	err := ValidateValueID(id)
	if err != nil {
		return fmt.Errorf("service: error to update user, validation failed")
	}
	return s.rps.Update(ctx, id, person)
}
func (s *Service) SelectAllUsers(ctx context.Context) ([]*model.Person, error) { //get all users from DB without passwords and tokens
	return s.rps.SelectAll(ctx)
}
func (s *Service) DeleteUser(ctx context.Context, id string) error { //delete user from DB
	err := ValidateValueID(id)
	if err != nil {
		return fmt.Errorf("service: error to delete user, validation failed")
	}
	return s.rps.Delete(ctx, id)
}
func (s *Service) GetUserById(ctx context.Context, id string) (model.Person, error) { //get one user by id
	err := ValidateValueID(id)
	if err != nil {
		return model.Person{}, fmt.Errorf("service: error to get user by id, validation failed")
	}
	return s.rps.SelectById(ctx, id)
}

func ValidateValueID(id string) error {
	err := validate.Var(id, "required, min=36")
	if err != nil {
		return fmt.Errorf("id length couldnt be less then 36,~%v", err)
	}
	return nil
}
