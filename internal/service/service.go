package service

import (
	"awesomeProject/internal/cache"
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
	"fmt"
)

var JwtKey = []byte("super-key") //key fo generation and check tokens

type Service struct { //service new
	rps       repository.Repository
	userCache *cache.UserCache
}

func NewService(NewRps repository.Repository, userCache *cache.UserCache) *Service { //create
	return &Service{rps: NewRps, userCache: userCache}
}

func (s *Service) UpdateUser(ctx context.Context, id string, person *model.Person) error { //update user
	return s.rps.Update(ctx, id, person)
}
func (s *Service) SelectAllUsers(ctx context.Context) ([]*model.Person, error) { //get all users from DB without passwords and tokens
	users, err := s.rps.SelectAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to select all users from db, %e", err)
	}
	err = s.AddUsersToCache(users, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to add users into the cache, %e", err)
	}
	return users, nil
}
func (s *Service) DeleteUser(ctx context.Context, id string) error { //delete user from DB
	return s.rps.Delete(ctx, id)
}
func (s *Service) GetUserById(ctx context.Context, id string) (model.Person, error) { //get one user by id
	user, err := s.rps.SelectById(ctx, id)
	if err != nil {
		return model.Person{}, fmt.Errorf("failed to select all users from db, %e", err)
	}
	err = s.userCache.AddToCache(ctx, &user)
	if err != nil {
		return model.Person{}, fmt.Errorf("failed to add users into the cache, %e", err)
	}
	return user, nil
}
func (s *Service) GetUserFromCache(ctx context.Context) (model.Person, bool, error) {
	return s.userCache.GetUserByIdFromCache(ctx)
}
func (s *Service) GetAllUsersFromCache(ctx context.Context) ([]*model.Person, bool, error) {
	return s.userCache.GetAllUsersFromCache(ctx)
}
func (s *Service) DeleteUserFromCache(ctx context.Context) error {
	return s.userCache.DeleteUserFromCache(ctx)
}
func (s *Service) AddUsersToCache(person []*model.Person, ctx context.Context) error {
	return s.userCache.AddAllUsersToCache(person, ctx)
}
