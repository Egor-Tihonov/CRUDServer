package service

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	rps repository.Repository
}

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

func NewService(NewRps repository.Repository) *Service {
	return &Service{rps: NewRps}
}

func (s *Service) CreateUser(ctx context.Context, person *model.Person) (error, string) { //registration
	hPassword, err := HashPassword(person.Password)
	if err != nil {
		return err, ""
	}
	person.Password = hPassword
	err, newId := s.rps.Create(ctx, person)
	if err != nil {
		return fmt.Errorf("service: registration failed: %v", err), ""
	}
	return nil, newId
}

func (s *Service) UpdateUser(ctx context.Context, id string, person *model.Person) error {
	return s.rps.Update(ctx, id, person)
}
func (s *Service) SelectAllUsers(ctx context.Context) ([]*model.Person, error) {
	return s.rps.SelectAll(ctx)
}
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	return s.rps.Delete(ctx, id)
}
func (s *Service) GetUserById(ctx context.Context, id string, password string) (model.Person, error) { //authentication
	authUser := model.Person{}
	authUser, err := s.rps.SelectById(ctx, id)
	if err != nil {
		return authUser, fmt.Errorf("service: authentication failed - %v", err)
	}
	incoming := []byte(password)
	existing := []byte(authUser.Password)
	err = bcrypt.CompareHashAndPassword(existing, incoming)
	if err != nil {
		return model.Person{}, fmt.Errorf("incorrect password: %v", err)
	}
	authUser.Password = password
	return authUser, nil
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("service: password is equals to zero")
	}
	bytesPassword := []byte(password)
	hashedBytesPassword, err := bcrypt.GenerateFromPassword(bytesPassword, bcrypt.DefaultCost)
	/*a1, _ := bcrypt.GenerateFromPassword(bytesPassword, bcrypt.DefaultCost)
	err = bcrypt.CompareHashAndPassword(a1, []byte(password))*/
	//log.Error(err)
	if err != nil {
		return "", err
	}
	hashPassword := string(hashedBytesPassword[:])
	return hashPassword, nil
}
