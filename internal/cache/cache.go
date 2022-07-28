package cache

import (
	"awesomeProject/internal/model"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

type UserCache struct {
	redisClient *redis.Client
}

func NewCache(rdsClient *redis.Client) *UserCache {
	return &UserCache{redisClient: rdsClient}
}

func (u *UserCache) AddToCache(ctx context.Context, person *model.Person) error {
	user, err := json.Marshal(person)
	if err != nil {
		log.Errorf("cache: failed add user to cache, %e", err)
		return err
	}
	err = u.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "user",
		ID:     "0-*",
		Values: map[string]interface{}{"About": user},
	}).Err()
	if err != nil {
		log.Errorf("cache: failed add user to cache, %e", err)
		return err
	}
	return nil
}

func (u *UserCache) GetUserByIdFromCache(ctx context.Context) (model.Person, bool, error) {
	result, err := u.redisClient.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"user", "0"},
		Count:   1,
		Block:   1 * time.Millisecond,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return model.Person{}, false, nil
		}
		log.Errorf("failed get user by id from cache: %e", err)
		return model.Person{}, false, err
	}
	bytes := result[0].Messages[0]
	msg := bytes.Values
	msgString := msg["About"].(string)
	person := model.Person{}
	err = json.Unmarshal([]byte(msgString), &person)
	if err != nil {
		log.Errorf("failed get user by id from cache: %e", err)
		return model.Person{}, false, err
	}
	return person, true, nil
}
func (u *UserCache) DeleteUserFromCache(ctx context.Context) error {
	err := u.redisClient.FlushAll(ctx).Err()
	if err != nil {
		log.Errorf("failed to delete user from cache, %e", err)
		return err
	}
	return nil
}

func (u *UserCache) GetAllUsersFromCache(ctx context.Context) ([]*model.Person, bool, error) {
	result, err := u.redisClient.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"all-users", "0"},
		Count:   1,
		Block:   1 * time.Millisecond,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		log.Errorf("failed get user by id from cache: %e", err)
		return nil, false, err
	}
	bytes := result[0].Messages[0]
	msg := bytes.Values
	msgString := msg["About"].(string)
	var persons []*model.Person

	err = json.Unmarshal([]byte(msgString), &persons)
	if err != nil {
		log.Errorf("failed to unmarshal json, %e", err)
		return nil, false, err
	}

	return persons, true, nil
}
func (u *UserCache) AddAllUsersToCache(person []*model.Person, ctx context.Context) error {
	user, err := json.Marshal(person)
	if err != nil {
		log.Errorf("cache: failed add all users to cache, %e", err)
		return err
	}
	err = u.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "all-users",
		ID:     "0-*",
		Values: map[string]interface{}{"About": user},
	}).Err()
	if err != nil {
		log.Errorf("cache: failed add all users to cache, %e", err)
		return err
	}
	return nil
}
