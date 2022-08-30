package domain

type Services interface {
	GetAccountService() AccountService
}
