package domain

type Repositories interface {
	GetAccountRepository() AccountRepository
}
