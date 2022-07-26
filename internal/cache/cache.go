package cache

import (
	"awesomeProject/internal/model"
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
)

type UserCache struct {
	redisClient *redis.Client
}

func NewCache(rdsClient *redis.Client) *UserCache {
	return &UserCache{redisClient: rdsClient}
}

func (userCache *UserCache) AddToCache(ctx context.Context, person model.Person) error {
	err := userCache.redisClient.Set(ctx, "user", "ID: "+person.ID+" Name: "+person.Name+" Age: "+strconv.Itoa(person.Age)+" Works: "+strconv.FormatBool(person.Works), 12*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (userCache *UserCache) GetUserByIdFromCache(ctx context.Context) (string, bool, error) {
	user, err := userCache.redisClient.Get(ctx, "user").Result()
	if err == redis.Nil {
		return "", false, nil
	}
	return user, true, nil
}

func (UserCache *UserCache) DeleteFromCache() {
	UserCache.redisClient.Del(context.Background(), "user")
}
