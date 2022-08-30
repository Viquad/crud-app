package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Viquad/crud-app/internal/domain"
	"github.com/golang-jwt/jwt/v4"
)

type UserService struct {
	repo struct {
		user  domain.UserRepository
		token domain.TokenRepository
	}
	hasher          PasswordHasher
	hmacSecret      []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(repos Repositories, hasher PasswordHasher, secret []byte, accessttl, refreshttl time.Duration) *UserService {
	return &UserService{
		repo: struct {
			user  domain.UserRepository
			token domain.TokenRepository
		}{
			user:  repos.GetUserRepository(),
			token: repos.GetTokenRepository(),
		},
		hasher:          hasher,
		hmacSecret:      secret,
		accessTokenTTL:  accessttl,
		refreshTokenTTL: refreshttl,
	}
}

func (s *UserService) Create(ctx context.Context, input domain.SignUpInput) error {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	input.Password = password

	return s.repo.user.Create(ctx, input)
}

func (s *UserService) GetTokenByCredentials(ctx context.Context, input domain.SignInInput) (string, string, error) {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return "", "", err
	}

	input.Password = password

	user, err := s.repo.user.GetByCredentials(ctx, input)
	if err != nil {
		return "", "", err
	}

	return s.generateTokens(ctx, user.Id)
}

func (s *UserService) GetByCredentials(ctx context.Context, input domain.SignInInput) (*domain.User, error) {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	input.Password = password

	return s.repo.user.GetByCredentials(ctx, input)
}

func (s *UserService) ParseToken(ctx context.Context, tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return 0, domain.ErrInvalidClaims
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return 0, domain.ErrAccessTokenExpired
	}

	id, err := strconv.ParseInt(claims.ID, 10, 64)
	if err != nil {
		return 0, domain.ErrInvalidId
	}

	return id, nil
}

func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := s.repo.token.Get(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return "", "", domain.ErrRefreshTokenExpired
	}

	return s.generateTokens(ctx, session.UserID)
}

func (s *UserService) generateTokens(ctx context.Context, userId int64) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        strconv.FormatInt(userId, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
	})

	accessToken, err := token.SignedString(s.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	refresh := domain.RefreshSession{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.refreshTokenTTL),
	}

	if err = s.repo.token.Create(ctx, refresh); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
