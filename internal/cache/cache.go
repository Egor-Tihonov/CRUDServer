package cache

import (
	"awesomeProject/internal/model"
	"context"
	"github.com/go-redis/redis/v9"
	"strconv"
	"time"
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

func (userCache *UserCache) GetUserByIdFromCache(ctx context.Context) (error, string, bool) {
	order, err := userCache.redisClient.Get(ctx, "user").Result()
	if err == redis.Nil {
		return nil, "", false
	}
	return nil, order, true
}
