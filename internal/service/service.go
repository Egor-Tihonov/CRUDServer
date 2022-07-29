package service

import (
	"awesomeProject/internal/cache"
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
	"fmt"
)

var JwtKey = []byte("super-key") // key fo generation and check tokens

type Service struct { // Service new
	rps       repository.Repository
	userCache *cache.UserCache
}

func NewService(newRps repository.Repository, userCache *cache.UserCache) *Service { // create
	return &Service{rps: newRps, userCache: userCache}
}

func (s *Service) UpdateUser(ctx context.Context, id string, person *model.Person) error { // update user
	err := s.rps.Update(ctx, id, person)
	if err != nil {
		return fmt.Errorf("failed to update users, %e", err)
	}
	err = s.userCache.AddToCache(ctx, person)
	if err != nil {
		return fmt.Errorf("failed to add update users into the cache, %e", err)
	}
	return nil
}
func (s *Service) SelectAllUsers(ctx context.Context) ([]*model.Person, error) { // get all users from DB without passwords and tokens
	users, found, err := s.userCache.GetAllUsersFromCache(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to select all users from db, %e", err)
	}
	if !found {
		users, err = s.rps.SelectAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to select all users from db, %e", err)
		}
		err = s.userCache.AddAllUsersToCache(users, ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to add users into the cache, %e", err)
		}
		return users, nil
	}
	return users, nil
}
func (s *Service) DeleteUser(ctx context.Context, id string) error { // delete user from DB
	_, found, err := s.userCache.GetUserByIDFromCache(ctx)
	if err != nil {
		return err
	}
	if !found {
		return s.rps.Delete(ctx, id)
	}
	err = s.userCache.DeleteUserFromCache(ctx)
	if err != nil {
		return fmt.Errorf("service: error while deleting user from cache, %e", err)
	}
	err = s.rps.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service: error while deleting user from DB, %e", err)
	}
	return nil
}
func (s *Service) GetUserByID(ctx context.Context, id string) (model.Person, error) { // get one user by id
	user, found, err := s.userCache.GetUserByIDFromCache(ctx)
	if err != nil {
		return model.Person{}, fmt.Errorf("failed to select user from cache, %e", err)
	}
	if !found {
		user, err = s.rps.SelectByID(ctx, id)
		if err != nil {
			return model.Person{}, fmt.Errorf("failed to select user from cache, %e", err)
		}
		err = s.userCache.AddToCache(ctx, &user)
		if err != nil {
			return model.Person{}, fmt.Errorf("failed to select user from cache, %e", err)
		}
		return user, nil
	}
	return user, nil
}

func (s *Service) DeleteFromCache(ctx context.Context) error {
	_, found, err := s.userCache.GetUserByIDFromCache(ctx)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	err = s.userCache.DeleteUserFromCache(ctx)
	if err != nil {
		return fmt.Errorf("service: error while deleting user from cache, %e", err)
	}
	return nil
}
