package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Viquad/crud-app/internal/domain"
	cache "github.com/Viquad/simple-cache"
	"github.com/sirupsen/logrus"
)

const cache_key_template = "user[%d]/account[%d]"
const listId int64 = 0

type AccountService struct {
	repo  domain.AccountRepository
	cache cache.Cache
	ttl   time.Duration
}

func NewAccountService(repo domain.AccountRepository, cache cache.Cache, ttl time.Duration) *AccountService {
	return &AccountService{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (s *AccountService) Create(ctx context.Context, input domain.AccountCreateInput) (*domain.Account, error) {
	userId, ok := ctx.Value(domain.UserIdKey).(int64)
	if !ok {
		return nil, domain.ErrInvalidId
	}

	account, err := s.repo.Create(ctx, input)
	if err == nil {
		s.cache.Set(cacheKey(userId, account.Id), account, s.ttl)
		s.cache.Delete(cacheKey(userId, listId))
	}

	return account, err
}

func (s *AccountService) GetById(ctx context.Context, id int64) (account *domain.Account, err error) {
	userId, ok := ctx.Value(domain.UserIdKey).(int64)
	if !ok {
		return nil, domain.ErrInvalidId
	}

	i, err := s.cache.Get(cacheKey(userId, id))
	if err == nil {
		logrus.WithFields(logrus.Fields{
			"context": "AccountService.GetById()",
		}).Debug("Get account from cache")
		account, ok = i.(*domain.Account)
	}

	if !ok {
		logrus.WithFields(logrus.Fields{
			"context": "AccountService.GetById()",
		}).Debug("Get account from repo")
		account, err = s.repo.GetById(ctx, id)
	}

	if err == nil {
		s.cache.Set(cacheKey(userId, id), account, s.ttl)
	}

	return account, err
}

func (s *AccountService) List(ctx context.Context) (accounts []domain.Account, err error) {
	userId, ok := ctx.Value(domain.UserIdKey).(int64)
	if !ok {
		return nil, domain.ErrInvalidId
	}

	i, err := s.cache.Get(cacheKey(userId, listId))
	if err == nil {
		logrus.WithFields(logrus.Fields{
			"context": "AccountService.List()",
		}).Debug("Get accounts from cache")
		accounts, ok = i.([]domain.Account)
	}

	if !ok {
		logrus.WithFields(logrus.Fields{
			"context": "AccountService.List()",
		}).Debug("Get accounts from repo")
		accounts, err = s.repo.List(ctx)
	}

	if err == nil {
		s.cache.Set(cacheKey(userId, listId), accounts, s.ttl)
	}

	return accounts, err
}

func (s *AccountService) UpdateById(ctx context.Context, id int64, inp domain.AccountUpdateInput) (*domain.Account, error) {
	userId, ok := ctx.Value(domain.UserIdKey).(int64)
	if !ok {
		return nil, domain.ErrInvalidId
	}

	account, err := s.repo.UpdateById(ctx, id, inp)
	if err == nil {
		s.cache.Set(cacheKey(userId, account.Id), account, s.ttl)
		s.cache.Delete(cacheKey(userId, listId))
	}

	return s.repo.UpdateById(ctx, id, inp)
}

func (s *AccountService) DeleteById(ctx context.Context, id int64) error {
	userId, ok := ctx.Value(domain.UserIdKey).(int64)
	if !ok {
		return domain.ErrInvalidId
	}

	err := s.repo.DeleteById(ctx, id)
	if err == nil {
		s.cache.Delete(cacheKey(userId, id))
	}

	return err
}

func cacheKey(user_id, id int64) string {
	return fmt.Sprintf(cache_key_template, user_id, id)
}
