package service

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var validate = validator.New()

func (s *Service) Authentication(ctx context.Context, id string, password string) (string, string, error) { //authentication
	err := ValidateValue(password)
	if err != nil {
		return "", "", fmt.Errorf("service: error with auth")
	}
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

func (s *Service) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, error) { //refresh our tokens
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
	if err != nil {
		return "", "", fmt.Errorf("service: can't generate refresh token - %v", err)
	}
	err = rps.UpdateAuth(ctx, person.ID, refreshTokenStr) //add into user refresh token
	if err != nil {
		return "", "", fmt.Errorf("service: can't generate refresh token - %v", err)
	}
	return accessTokenStr, refreshTokenStr, nil
}

func (s *Service) UpdateUserAuth(ctx context.Context, id string, refreshToken string) error { //add into user refresh token
	return s.rps.UpdateAuth(ctx, id, refreshToken)
}

func (s *Service) Registration(ctx context.Context, person *model.Person) (string, error) { //users`s registration
	err := ValidateStruct(person)
	if err != nil {
		return "", fmt.Errorf("service: error with creating new user,%v", err)
	}
	hPassword, err := HashPassword(person.Password) //check password (authentication)
	if err != nil {
		return " ", err
	}
	person.Password = hPassword
	newId, err := s.rps.Create(ctx, person)
	if err != nil {
		return "", fmt.Errorf("service: registration failed: %v", err)
	}
	return newId, nil
}

func HashPassword(password string) (string, error) {
	bytesPassword := []byte(password)
	hashedBytesPassword, err := bcrypt.GenerateFromPassword(bytesPassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashPassword := string(hashedBytesPassword[:])
	return hashPassword, nil
}

func ValidateStruct(person *model.Person) error {
	err := validate.Struct(person)
	if err != nil {
		return fmt.Errorf("error with validate user, check your name(min length = 6),password(min length = 8) and age couldnt be less then 0 or greater than 200,~ %v", err)
	}
	return nil
}

func ValidateValue(password string) error {
	err := validate.Var(password, "required,min=8")
	if err != nil {
		return fmt.Errorf("password length couldnt be less then 8,~%v", err)
	}
	return nil
}
