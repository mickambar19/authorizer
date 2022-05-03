package usecase

import (
	"time"

	"github.com/mickambar19/authorizer/src/helpers"
	"github.com/mickambar19/authorizer/src/model"
)

type TransactionService interface {
	AddTransaction(int, string, time.Time)
	GetTransactionFromOffset(int) (bool, model.Transaction)
	GetTransactionsBy(string, int, time.Time, int) []model.Transaction
}

type AccountService interface {
	AccountEnabled() bool
	AccountCardEnabled() bool
	GetAccount() model.Account
	UpdateAvailableLimit(int) model.Account
}

type Transaction struct {
	ts TransactionService
	as AccountService
}

func NewTransaction(ts TransactionService, as AccountService) *Transaction {
	return &Transaction{
		ts,
		as,
	}
}

func (t *Transaction) CreateTransaction(amount int, merchant string, date time.Time) (model.Account, []model.Violation) {
	account := t.as.GetAccount()
	violations := []model.Violation{}

	if !t.as.AccountEnabled() {
		return account, []model.Violation{"account-not-initialized"}
	}

	if !t.as.AccountCardEnabled() {
		return account, []model.Violation{"card-not-active"}
	}

	if !helpers.WithinLimit(account.AvailableLimit, amount) {
		violations = append(violations, "insufficient-limit")
	}

	intervalMinutes := 2
	exists, lastThirdTransaction := t.ts.GetTransactionFromOffset(3)

	if exists && helpers.AreDatesWithinInterval(lastThirdTransaction.Time, date, intervalMinutes) {
		violations = append(violations, "high-frequency-small-interval")
	}

	dupedTransactions := t.ts.GetTransactionsBy(merchant, amount, date, intervalMinutes)
	if len(dupedTransactions) > 0 {
		violations = append(violations, "doubled-transaction")
	}

	if len(violations) == 0 {
		t.ts.AddTransaction(amount, merchant, date)
		account = t.as.UpdateAvailableLimit(amount)
	}

	return account, violations
}
