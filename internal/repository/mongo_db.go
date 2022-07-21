package repository

import (
	"awesomeProject/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MRepository struct { //mongo
	Pool *mongo.Client
}

func (m *MRepository) Create(ctx context.Context, person *model.Person) (string, error) {
	if person.Age < 0 || person.Age > 180 {
		return "", error(fmt.Errorf("mongo repository: error with create, age must be more then 0 and less then 180"))
	}
	newId := uuid.New().String()
	collection := m.Pool.Database("person").Collection("person")
	_, err := collection.InsertOne(ctx, bson.D{
		{Key: "id", Value: newId},
		{Key: "name", Value: person.Name},
		{Key: "works", Value: person.Works},
		{Key: "age", Value: person.Age},
		{Key: "password", Value: person.Password},
		{Key: "refreshtoken", Value: person.RefreshToken},
	})
	if err != nil {
		return "", fmt.Errorf("mongo: unable to create new user: %v", err)
	}
	return newId, nil
}

func (m *MRepository) Update(ctx context.Context, id string, person *model.Person) error {
	if person.Age < 0 || person.Age > 180 {
		return error(fmt.Errorf("mongo repository: error with create, age must be more then 0 and less then 180"))
	}
	collection := m.Pool.Database("person").Collection("person")
	_, err := collection.UpdateOne(ctx, bson.D{{"id", id}}, bson.D{{"$set", bson.D{
		{"name", person.Name},
		{"works", person.Works},
		{"age", person.Age},
	}}})
	if err != nil {
		return fmt.Errorf("mongo: unable to update user %v", err)
	}
	return nil
}
func (m *MRepository) UpdateAuth(ctx context.Context, id string, refreshToken string) error {
	collection := m.Pool.Database("person").Collection("person")
	_, err := collection.UpdateOne(ctx, bson.D{{"id", id}}, bson.D{{"$set", bson.D{
		{"refreshtoken", refreshToken},
	}}})
	if err != nil {
		return fmt.Errorf("mongo: unable to update user %v", err)
	}
	return nil
}

func (m *MRepository) SelectAll(ctx context.Context) ([]*model.Person, error) {
	var users []*model.Person
	collection := m.Pool.Database("person").Collection("person")
	c, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("mongo: unable to select all users %v", err)
	}
	for c.Next(ctx) {
		user := model.Person{}
		err := c.Decode(&user)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}
	return users, nil

}

func (m *MRepository) Delete(ctx context.Context, id string) error {
	collection := m.Pool.Database("person").Collection("person")
	_, err := collection.DeleteOne(ctx, bson.D{{"id", id}})
	if err != nil {
		return fmt.Errorf("mongo: unable to delete user, %v", err)
	}
	return nil
}

func (m *MRepository) SelectById(ctx context.Context, id string) (model.Person, error) {
	user := model.Person{}
	collection := m.Pool.Database("person").Collection("person")
	err := collection.FindOne(ctx, bson.D{{"id", id}}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
func (m *MRepository) SelectByIdAuth(ctx context.Context, id string) (model.Person, error) {
	user := model.Person{}
	collection := m.Pool.Database("person").Collection("person")
	err := collection.FindOne(ctx, bson.D{{"id", id}, {"name", 0}, {"works", 0}, {"age", 0}, {"password", 0}}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
