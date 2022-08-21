package service

import (
	"context"
	"strconv"
	"time"

	"github.com/Viquad/crud-app/internal/domain"
	cache "github.com/Viquad/simple-cache"
	"github.com/sirupsen/logrus"
)

var cache_salt = "_account"
var listId int64 = 0

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

func (s *AccountService) Create(ctx context.Context, accountInput domain.Account) (*domain.Account, error) {
	account, err := s.repo.Create(ctx, accountInput)
	if err == nil {
		s.cache.Set(idToString(account.Id), account, s.ttl)
	}

	return account, err
}

func (s *AccountService) GetById(ctx context.Context, id int64) (account *domain.Account, err error) {
	i, err := s.cache.Get(idToString(id))
	var ok bool
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
		s.cache.Set(idToString(id), account, s.ttl)
	}

	return account, err
}

func (s *AccountService) List(ctx context.Context) (accounts []domain.Account, err error) {
	i, err := s.cache.Get(idToString(listId))
	var ok bool
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
		s.cache.Set(idToString(listId), accounts, s.ttl)
	}

	return accounts, err
}

func (s *AccountService) UpdateById(ctx context.Context, id int64, inp domain.AccountUpdateInput) (*domain.Account, error) {
	account, err := s.repo.UpdateById(ctx, id, inp)
	if err == nil {
		s.cache.Set(idToString(account.Id), account, s.ttl)
	}

	return s.repo.UpdateById(ctx, id, inp)
}

func (s *AccountService) DeleteById(ctx context.Context, id int64) error {
	err := s.repo.DeleteById(ctx, id)
	if err == nil {
		s.cache.Delete(idToString(id))
	}

	return err
}

func idToString(id int64) string {
	return strconv.FormatInt(id, 10) + cache_salt
}
