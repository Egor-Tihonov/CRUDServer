package cache

import (
	"awesomeProject/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
	"time"
)

type UserCache struct {
	redisClient *redis.Client
}

func NewCache(rdsClient *redis.Client) *UserCache {
	return &UserCache{redisClient: rdsClient}
}

func (userCache *UserCache) AddToCache(ctx context.Context, person *model.Person) error {
	user, err := json.Marshal(person)
	if err != nil {
		return fmt.Errorf("cache: failed add user to cache, %e", err)
	}
	err = userCache.redisClient.Set(ctx, "user", user, 12*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("cache: failed add user to cache, %e", err)
	}
	return nil
}

func (userCache *UserCache) GetUserByIdFromCache(ctx context.Context) (string, bool, error) {
	user, err := userCache.redisClient.Get(ctx, "user").Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, fmt.Errorf("cache: failed get user from cache, %e", err)
	}

	return user, true, nil
}

func (userCache *UserCache) DeleteUsersFromCache(ctx context.Context) error {
	err := userCache.redisClient.Del(ctx, "allusers").Err()
	if err != nil {
		return fmt.Errorf("cache: failed delete user from cache, %e", err)
	}
	return nil
}

func (userCache *UserCache) GetAllUsersFromCache(ctx context.Context) (string, bool, error) {
	users, err := userCache.redisClient.Get(ctx, "allusers").Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, fmt.Errorf("cache: failed get users from cache, %e", err)
	}
	return users, true, nil
}
func (userCache *UserCache) AddAllUsersToCache(person []*model.Person, ctx context.Context) error {
	users, err := json.Marshal(person)
	if err != nil {
		return fmt.Errorf("cache: failed add user to cache, %e", err)
	}
	err = userCache.redisClient.Set(ctx, "allusers", users, 1*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}
