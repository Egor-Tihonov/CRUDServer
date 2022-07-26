package service

import (
	"awesomeProject/internal/cache"
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
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
	return s.rps.SelectAll(ctx)
}
func (s *Service) DeleteUser(ctx context.Context, id string) error { //delete user from DB
	return s.rps.Delete(ctx, id)
}
func (s *Service) GetUserById(ctx context.Context, id string) (model.Person, error) { //get one user by id
	return s.rps.SelectById(ctx, id)
}
func (s *Service) GetUserByIdFromCache(ctx context.Context) (string, bool, error) {
	return s.userCache.GetUserByIdFromCache(ctx)
}
func (s *Service) AddToCache(ctx context.Context, person model.Person) error {
	return s.userCache.AddToCache(ctx, person)
}
func (s *Service) DeleteFromCache() {
	s.userCache.DeleteFromCache()
}
