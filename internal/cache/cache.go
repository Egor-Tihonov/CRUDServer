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
	err = userCache.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "user",
		ID:     "0-*",
		Values: map[string]interface{}{"About": user},
	}).Err()
	if err != nil {
		return fmt.Errorf("cache: failed add user to cache, %e", err)
	}
	/*err = userCache.redisClient.Set(ctx, "user", user, 12*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("cache: failed add user to cache, %e", err)
	}*/
	return nil
}

func (userCache *UserCache) GetUserByIdFromCache(ctx context.Context) (model.Person, bool, error) {
	result, err := userCache.redisClient.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"user", "0"},
		Count:   1,
		Block:   1 * time.Millisecond,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return model.Person{}, false, nil
		}
		return model.Person{}, false, fmt.Errorf("failed get user by id from cache: %e", err)
	}
	bytes := result[0].Messages[0]
	msg := bytes.Values
	msgString := msg["About"].(string)
	person := model.Person{}
	err = json.Unmarshal([]byte(msgString), &person)
	if err != nil {
		fmt.Print(err)
	}
	return person, true, nil
}
func (userCache *UserCache) DeleteUserFromCache(ctx context.Context) error {
	/*err := userCache.redisClient.Del(ctx, "allusers").Err()
	if err != nil {
		return fmt.Errorf("cache: failed delete user from cache, %e", err)
	}
	return nil*/
	err := userCache.redisClient.FlushAll(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user from cache")
	}
	return nil
}

func (userCache *UserCache) GetAllUsersFromCache(ctx context.Context) ([]*model.Person, bool, error) {
	/*users, err := userCache.redisClient.Get(ctx, "allusers").Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, fmt.Errorf("cache: failed get users from cache, %e", err)
	}
	return users, true, nil*/
	result, err := userCache.redisClient.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"all-users", "0"},
		Count:   1,
		Block:   1 * time.Millisecond,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed get user by id from cache: %e", err)
	}
	bytes := result[0].Messages[0]
	msg := bytes.Values
	msgString := msg["About"].(string)
	var persons []*model.Person

	err = json.Unmarshal([]byte(msgString), &persons)
	if err != nil {
		fmt.Print(err)
	}

	return persons, true, nil
}
func (userCache *UserCache) AddAllUsersToCache(person []*model.Person, ctx context.Context) error {
	/*
		users, err := json.Marshal(person)
		if err != nil {
			return fmt.Errorf("cache: failed add user to cache, %e", err)
		}
		err = userCache.redisClient.Set(ctx, "all-users", users, 1*time.Minute).Err()
		if err != nil {
			return err
		}
		return nil*/

	user, err := json.Marshal(person)
	if err != nil {
		return fmt.Errorf("cache: failed add user to cache, %e", err)
	}
	err = userCache.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "all-users",
		ID:     "0-*",
		Values: map[string]interface{}{"About": user},
	}).Err()
	if err != nil {
		return fmt.Errorf("cache: failed add user to cache, %e", err)
	}
	return nil
}
