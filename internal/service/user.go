package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Viquad/crud-app/internal/domain"
	"github.com/Viquad/crud-audit-service/pkg/domain/audit"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const auth_session = "auth_session"

type UserService struct {
	repo struct {
		user  domain.UserRepository
		token domain.TokenRepository
	}
	sessionsStore   sessions.Store
	hasher          PasswordHasher
	hmacSecret      []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	audit           audit.AuditServiceClient
}

func NewUserService(repos Repositories, store sessions.Store, audit audit.AuditServiceClient, hasher PasswordHasher, secret []byte, accessttl, refreshttl time.Duration) *UserService {
	return &UserService{
		repo: struct {
			user  domain.UserRepository
			token domain.TokenRepository
		}{
			user:  repos.GetUserRepository(),
			token: repos.GetTokenRepository(),
		},
		sessionsStore:   store,
		hasher:          hasher,
		hmacSecret:      secret,
		accessTokenTTL:  accessttl,
		refreshTokenTTL: refreshttl,
		audit:           audit,
	}
}

func (s *UserService) Create(ctx context.Context, input domain.SignUpInput) error {
	password, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	input.Password = password

	user_id, err := s.repo.user.Create(ctx, input)
	if err != nil {
		return err
	}

	_, err = s.audit.Log(ctx, &audit.LogRequest{
		Action:    audit.LogRequest_REGISTER,
		UserId:    user_id,
		Timestamp: timestamppb.Now(),
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"context": "UserService.Create",
			"problem": "audit service Log()",
		}).Error(err.Error())
	}

	return nil
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

	access, refresh, err := s.generateTokens(ctx, user.Id)
	if err != nil {
		return "", "", err
	}

	_, err = s.audit.Log(ctx, &audit.LogRequest{
		Action:    audit.LogRequest_LOGIN,
		UserId:    user.Id,
		Timestamp: timestamppb.Now(),
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"context": "UserService.GetTokenByCredentials",
			"problem": "audit service Log()",
		}).Error(err.Error())
	}

	return access, refresh, nil
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

	access, refresh, err := s.generateTokens(ctx, session.UserID)
	if err != nil {
		return "", "", err
	}

	_, err = s.audit.Log(ctx, &audit.LogRequest{
		Action:    audit.LogRequest_REFRESH,
		UserId:    session.UserID,
		Timestamp: timestamppb.Now(),
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"context": "UserService.RefreshTokens",
			"problem": "audit service Log()",
		}).Error(err.Error())
	}

	return access, refresh, nil
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

func (s *UserService) InitSession(c *gin.Context, input domain.SignInInput) error {
	user, err := s.GetByCredentials(c.Request.Context(), input)
	if err != nil {
		return err
	}

	session, err := s.sessionsStore.New(c.Request, auth_session)
	if err != nil {
		return err
	}

	session.Values[string(domain.UserIdKey)] = user.Id

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		return err
	}

	_, err = s.audit.Log(c.Request.Context(), &audit.LogRequest{
		Action:    audit.LogRequest_LOGIN,
		UserId:    user.Id,
		Timestamp: timestamppb.Now(),
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"context": "UserService.RefreshTokens",
			"problem": "audit service Log()",
		}).Error(err.Error())
	}

	return nil
}

func (s *UserService) GetSession(c *gin.Context) (int64, error) {
	session, err := s.sessionsStore.Get(c.Request, auth_session)
	if err != nil {
		return 0, err
	}

	id, ok := session.Values[string(domain.UserIdKey)].(int64)
	if !ok {
		return 0, domain.ErrInvalidId
	}

	return id, nil
}

func (s *UserService) DropSession(c *gin.Context) error {
	session, _ := s.sessionsStore.Get(c.Request, auth_session)
	session.Options.MaxAge = -1

	err := session.Save(c.Request, c.Writer)
	if err != nil {
		return err
	}

	user_id, ok := session.Values[string(domain.UserIdKey)].(int64)
	if !ok {
		return nil
	}

	_, err = s.audit.Log(c.Request.Context(), &audit.LogRequest{
		Action:    audit.LogRequest_LOGOUT,
		UserId:    user_id,
		Timestamp: timestamppb.Now(),
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"context": "UserService.RefreshTokens",
			"problem": "audit service Log()",
		}).Error(err.Error())
	}

	return nil
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
