package service

import (
	todo "Todo_rest_api"
	"Todo_rest_api/pkg/repository"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTl   = time.Hour * time.Duration(1)
	signingKey = "asd234asd"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}
func (s *AuthService) CreateUser(user todo.User) (int, error) {
	hash, err := generatePasswordHash(user.Password)
	if err != nil {
		return 0, fmt.Errorf("invalid password: %w", err)
	}
	user.Password = hash
	return s.repo.CreateUser(user)

}

type TokenClaims struct {
	jwt.RegisteredClaims
	UserId int `json:"user_id"`
}

func (s *AuthService) GenerateToken(username string, password string) (string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {

	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	return claims.UserId, nil
}

func generatePasswordHash(password string) (string, error) {
	if strings.TrimSpace(password) == "" {
		return "", errors.New("password is empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}
