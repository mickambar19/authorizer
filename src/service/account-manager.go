package service

import (
	"github.com/mickambar19/authorizer/src/model"
)

type AccountManager struct {
	account *model.Account
}

func NewAccountManager() *AccountManager {
	return &AccountManager{}
}

func (a *AccountManager) AccountEnabled() bool {
	return a.account != nil
}

func (a *AccountManager) AccountCardEnabled() bool {
	if a.account == nil {
		return false
	}
	return a.account.ActiveCard
}

func (a *AccountManager) CreateAccount(activeCard bool, availableLimit int) model.Account {
	a.account = &model.Account{
		ActiveCard:     activeCard,
		AvailableLimit: availableLimit,
	}
	return *a.account
}

func (a *AccountManager) GetAccount() model.Account {
	if a.account == nil {
		return model.Account{}
	}
	return *a.account
}

func (a *AccountManager) UpdateAvailableLimit(amountSpent int) model.Account {
	a.account.AvailableLimit -= amountSpent
	return *a.account
}
