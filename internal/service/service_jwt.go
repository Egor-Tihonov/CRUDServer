package service

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	accessTokenWorkTime  = time.Now().Add(time.Minute * 5).Unix()
	refreshTokenWorkTime = time.Now().Add(time.Hour * 3).Unix()
)

func (s *Service) Authentication(ctx context.Context, id, password string) (accessTokenStr, refreshTokenStr string, err error) { // Authentication
	authUser, err := s.rps.SelectByID(ctx, id)
	if err != nil {
		return "", "", fmt.Errorf("service: authentication failed - %v", err)
	}
	incoming := []byte(password)
	existing := []byte(authUser.Password)
	err = bcrypt.CompareHashAndPassword(existing, incoming) // check passwords
	if err != nil {
		return "", "", fmt.Errorf("incorrect password: %v", err)
	}
	authUser.Password = password

	return s.CreateJWT(s.rps, &authUser, ctx)
}

func (s *Service) RefreshToken(ctx context.Context, refreshTokenString string) (accessTokenStr, refreshTokenStr string, err error) { // refresh our tokens
	refreshToken, err := jwt.Parse(refreshTokenString, func(t *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	}) // parse it into string format
	if err != nil {
		log.Errorf("service: can't parse refresh token - %e", err)
		return "", "", err
	}
	if !refreshToken.Valid {
		return "", "", fmt.Errorf("service: expired refresh token")
	}
	claims := refreshToken.Claims.(jwt.MapClaims)
	userUUID := claims["jti"]
	if userUUID == "" {
		return "", "", fmt.Errorf("service: error while parsing claims, ID couldnt be empty")
	}
	person, err := s.rps.SelectByIDAuth(ctx, userUUID.(string))
	if err != nil {
		return "", "", fmt.Errorf("service: token refresh failed - %e", err)
	}
	if refreshTokenString != person.RefreshToken {
		return "", "", fmt.Errorf("service: invalid refresh token")
	}
	return s.CreateJWT(s.rps, &person, ctx)
}

func (s *Service) CreateJWT(rps repository.Repository, person *model.Person, ctx context.Context) (accessTokenStr, refreshTokenStr string, err error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)         // encrypt access token by SigningMethodHS256 method
	claimsA := accessToken.Claims.(jwt.MapClaims)          // fill access-token`s claims
	claimsA["exp"] = accessTokenWorkTime                   // work time
	claimsA["username"] = person.Name                      // payload
	accessTokenStr, err = accessToken.SignedString(JwtKey) // convert token to string format
	if err != nil {
		log.Errorf("service: can't generate access token - %v", err)
		return "", "", err
	}
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claimsR := refreshToken.Claims.(jwt.MapClaims)
	claimsR["username"] = person.Name
	claimsR["exp"] = refreshTokenWorkTime
	claimsR["jti"] = person.ID
	refreshTokenStr, err = refreshToken.SignedString(JwtKey)
	if err != nil {
		log.Errorf("service: can't generate access token - %v", err)
		return "", "", err
	}
	err = rps.UpdateAuth(ctx, person.ID, refreshTokenStr) // add into user refresh token
	if err != nil {
		log.Errorf("service: can't generate access token - %v", err)
		return "", "", err
	}
	return
}

func (s *Service) UpdateUserAuth(ctx context.Context, id, refreshToken string) error { // add into user refresh token
	return s.rps.UpdateAuth(ctx, id, refreshToken)
}

func (s *Service) Registration(ctx context.Context, person *model.Person) (string, error) { // users`s registration
	hPassword, err := HashPassword(person.Password)
	if err != nil {
		return " ", err
	}
	person.Password = hPassword
	newID, err := s.rps.Create(ctx, person)
	if err != nil {
		return "", err
	}

	return newID, nil
}

func HashPassword(password string) (string, error) {
	bytesPassword := []byte(password)
	hashedBytesPassword, err := bcrypt.GenerateFromPassword(bytesPassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashPassword := string(hashedBytesPassword)
	return hashPassword, nil
}
