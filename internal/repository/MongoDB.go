package repository

import (
	"awesomeProject/internal/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MRepository struct {
	Pool *mongo.Client
}

func (m MRepository) Create(ctx context.Context, person *model.Person) error {
	i := 1
	collection := m.Pool.Database("test").Collection("person")
	_, err := collection.InsertOne(ctx, bson.D{
		{Key: "id", Value: i},
		{Key: "name", Value: person.Name},
		{Key: "works", Value: person.Works},
		{Key: "age", Value: person.Age},
	})
	if err != nil {
		return fmt.Errorf("mongo: unable to create new user: %v", err)
	}
	i++
	return nil
}

func (m MRepository) Update(ctx context.Context, id int, person *model.Person) error {
	//TODO implement me
	panic("implement me")
}

func (m MRepository) SelectAll(ctx context.Context) ([]*model.Person, error) {
	//TODO implement me
	panic("implement me")
}

func (m MRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (m MRepository) SelectById(ctx context.Context, id int) (model.Person, error) {
	//TODO implement me
	panic("implement me")
}

func NewMongoConn(_Pool *mongo.Client) MRepository {
	return MRepository{Pool: _Pool}
}
