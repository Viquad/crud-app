package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Viquad/crud-app/internal/domain"
	"github.com/golang-jwt/jwt/v4"
)

type UserService struct {
	repo   domain.UserRepository
	hasher PasswordHasher

	hmacSecret []byte
	tokenTTL   time.Duration
}

func NewUserService(repo domain.UserRepository, hasher PasswordHasher, secret []byte, ttl time.Duration) *UserService {
	return &UserService{
		repo:       repo,
		hasher:     hasher,
		hmacSecret: secret,
		tokenTTL:   ttl,
	}
}

func (s *UserService) Create(ctx context.Context, input domain.SignUpInput) error {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	input.Password = password

	return s.repo.Create(ctx, input)
}

func (s *UserService) GetTokenByCredentials(ctx context.Context, input domain.SignInInput) (string, error) {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", err
	}

	input.Password = password

	user, err := s.repo.GetByCredentials(ctx, input)
	if err != nil {
		return "", err
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        strconv.FormatInt(user.Id, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(s.hmacSecret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (s *UserService) GetByCredentials(ctx context.Context, input domain.SignInInput) (*domain.User, error) {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	input.Password = password

	return s.repo.GetByCredentials(ctx, input)
}

func (s *UserService) ParseToken(ctx context.Context, tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.RegisteredClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	id, err := strconv.ParseInt(claims.ID, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id")
	}

	return id, nil
}
