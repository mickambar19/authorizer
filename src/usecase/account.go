package usecase

import "github.com/mickambar19/authorizer/src/model"

type AccountManager interface {
	AccountEnabled() bool
	CreateAccount(activeCard bool, availableImit int) model.Account
	GetAccount() model.Account
}

type Account struct {
	am AccountManager
}

func NewAccount(am AccountManager) *Account {
	return &Account{
		am,
	}
}

func (a *Account) CreateAccount(activeCard bool, availableLimit int) (model.Account, []model.Violation) {
	violations := []model.Violation{}

	if a.am.AccountEnabled() {
		return a.am.GetAccount(), append(violations, model.Violation("account-already-initialized"))
	}

	account := a.am.CreateAccount(activeCard, availableLimit)

	return account, violations
}
