package service

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var JwtKey = []byte("super-key") //key fo generation and check tokens

type Service struct {
	rps repository.Repository
}

func NewService(NewRps repository.Repository) *Service {
	return &Service{rps: NewRps}
}

func (s *Service) Registration(ctx context.Context, person *model.Person) (error, string) { //users`s registration
	hPassword, err := HashPassword(person.Password) //check password (authentication)
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

func (s *Service) UpdateUserAuth(ctx context.Context, id string, refreshToken string) error { //add into user refresh token
	return s.rps.UpdateAuth(ctx, id, refreshToken)
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
	user, err := s.rps.SelectById(ctx, id)
	if err != nil {
		return model.Person{}, fmt.Errorf("serivce: error with select user by id: %v", err)
	}
	return user, nil
}

func (s *Service) Authentication(ctx context.Context, id string, password string) (string, string, error) { //authentication
	authUser, err := s.rps.SelectById(ctx, id)
	if err != nil {
		return "", "", fmt.Errorf("service: authentication failed - %v", err)
	}
	incoming := []byte(password)
	existing := []byte(authUser.Password)
	err = bcrypt.CompareHashAndPassword(existing, incoming) //check passwords
	if err != nil {
		return "", "", fmt.Errorf("incorrect password: %v", err)
	}
	authUser.Password = password

	return CreateJWT(s.rps, &authUser, ctx)
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("service: password is equals to zero")
	}
	bytesPassword := []byte(password)
	hashedBytesPassword, err := bcrypt.GenerateFromPassword(bytesPassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashPassword := string(hashedBytesPassword[:])
	return hashPassword, nil
}

func CreateJWT(rps repository.Repository, person *model.Person, ctx context.Context) (string, string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)           //encrypt access token by SigningMethodHS256 method
	claimsA := accessToken.Claims.(jwt.MapClaims)            //fill access-token`s claims
	claimsA["exp"] = time.Now().Add(15 * time.Minute).Unix() //work time
	claimsA["username"] = person.Name                        //payload
	accessTokenStr, err := accessToken.SignedString(JwtKey)  //convert token to string format
	if err != nil {
		return "", "", fmt.Errorf("service: can't generate access token - %v", err)
	}
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claimsR := refreshToken.Claims.(jwt.MapClaims)
	claimsR["username"] = person.Name
	claimsR["exp"] = time.Now().Add(time.Hour * 3).Unix()
	claimsR["jti"] = person.ID
	refreshTokenStr, err := refreshToken.SignedString(JwtKey)
	err = rps.UpdateAuth(ctx, person.ID, refreshTokenStr) //add into user refresh token
	if err != nil {
		return "", "", fmt.Errorf("service: can't generate refresh token - %v", err)
	}
	return accessTokenStr, refreshTokenStr, nil
}

func (s Service) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, error) { //refresh our tokens
	refreshToken, err := jwt.Parse(refreshTokenString, func(t *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	}) //parse it into string format
	if err != nil {
		return "", "", fmt.Errorf("service: can't parse refresh token - %w", err)
	}
	if !refreshToken.Valid {
		return "", "", fmt.Errorf("service: expired refresh token")
	}
	claims := refreshToken.Claims.(jwt.MapClaims)
	userUUID := claims["jti"]
	if userUUID == "" {
		return "", "", fmt.Errorf("service: error while parsing claims")
	}
	person, err := s.rps.SelectByIdAuth(ctx, userUUID.(string))
	if err != nil {
		return "", "", fmt.Errorf("service: token refresh failed - %w", err)
	}
	if refreshTokenString != person.RefreshToken {
		return "", "", fmt.Errorf("service: invalid refresh token")
	}
	return CreateJWT(s.rps, &person, ctx)
}
