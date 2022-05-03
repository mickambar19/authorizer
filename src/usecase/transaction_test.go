package usecase

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mocks "github.com/mickambar19/authorizer/src/mock"
	"github.com/mickambar19/authorizer/src/model"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -source=transaction.go -destination=../mock/usecase_transaction.go -package=mocks

func TestCreateTransaction(t *testing.T) {
	assert := assert.New(t)
	now := time.Now()
	tests := []struct {
		name                    string
		shouldAddTransaction    bool
		amount                  int
		merchantName            string
		date                    time.Time
		hasAccountEnabled       bool
		hasCardEnabled          bool
		accountAvailableLimit   int
		lastThirdTransaction    model.Transaction
		hasLastThirdTransaction bool
		dupedTransactions       []model.Transaction
		expectedAccount         model.Account
		expectedViolations      []model.Violation
	}{
		{
			name:                  "Processing a transaction successfully",
			shouldAddTransaction:  true,
			amount:                80,
			date:                  now,
			hasAccountEnabled:     true,
			hasCardEnabled:        true,
			accountAvailableLimit: 100,
			lastThirdTransaction:  model.Transaction{},
			dupedTransactions:     []model.Transaction{},
			expectedAccount: model.Account{
				ActiveCard:     true,
				AvailableLimit: 20,
			},
			expectedViolations: []model.Violation{},
		},
		{
			name:                  "Transaction raised account not initialized",
			amount:                200,
			date:                  now,
			hasAccountEnabled:     false,
			hasCardEnabled:        false,
			accountAvailableLimit: 0,
			lastThirdTransaction:  model.Transaction{},
			dupedTransactions:     []model.Transaction{},
			expectedAccount:       model.Account{},
			expectedViolations: []model.Violation{
				"account-not-initialized",
			},
		},
		{
			name:                  "Transaction raised card not active",
			amount:                200,
			date:                  now,
			hasAccountEnabled:     true,
			hasCardEnabled:        false,
			accountAvailableLimit: 300,
			lastThirdTransaction:  model.Transaction{},
			dupedTransactions:     []model.Transaction{},
			expectedAccount: model.Account{
				ActiveCard:     false,
				AvailableLimit: 300,
			},
			expectedViolations: []model.Violation{
				"card-not-active",
			},
		},
		{
			name:                  "Transaction raised insufficient-limit",
			amount:                300,
			date:                  now,
			hasAccountEnabled:     true,
			hasCardEnabled:        true,
			accountAvailableLimit: 200,
			lastThirdTransaction:  model.Transaction{},
			dupedTransactions:     []model.Transaction{},
			expectedAccount: model.Account{
				ActiveCard:     true,
				AvailableLimit: 200,
			},
			expectedViolations: []model.Violation{
				"insufficient-limit",
			},
		},
		{
			name:                    "Transaction raised high frequency small interval",
			amount:                  20,
			date:                    now,
			hasAccountEnabled:       true,
			hasCardEnabled:          true,
			accountAvailableLimit:   60,
			hasLastThirdTransaction: true,
			lastThirdTransaction: model.Transaction{
				Amount:   20,
				Merchant: "Burger King",
				Time:     now.Add(-time.Second * 30),
			},
			dupedTransactions: []model.Transaction{},
			expectedAccount: model.Account{
				ActiveCard:     true,
				AvailableLimit: 60,
			},
			expectedViolations: []model.Violation{
				"high-frequency-small-interval",
			},
		},
		{
			name:                  "Duped transaction",
			merchantName:          "Subway",
			amount:                200,
			date:                  now,
			hasAccountEnabled:     true,
			hasCardEnabled:        true,
			accountAvailableLimit: 500,
			lastThirdTransaction:  model.Transaction{},
			dupedTransactions: []model.Transaction{
				{
					Merchant: "Subway",
					Amount:   200,
					Time:     now,
				},
			},
			expectedAccount: model.Account{
				AvailableLimit: 500,
				ActiveCard:     true,
			},
			expectedViolations: []model.Violation{
				"doubled-transaction",
			},
		},
		{
			name:                    "Transaction violates multiple business rules",
			amount:                  20,
			date:                    now,
			hasAccountEnabled:       true,
			hasCardEnabled:          true,
			accountAvailableLimit:   10,
			hasLastThirdTransaction: true,
			lastThirdTransaction: model.Transaction{
				Amount:   20,
				Merchant: "Burger King",
				Time:     now.Add(-time.Minute * 1),
			},
			dupedTransactions: []model.Transaction{},
			expectedAccount: model.Account{
				AvailableLimit: 10,
				ActiveCard:     true,
			},
			expectedViolations: []model.Violation{
				"insufficient-limit",
				"high-frequency-small-interval",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			transactionService := mocks.NewMockTransactionService(ctrl)
			accountService := mocks.NewMockAccountService(ctrl)

			accountService.EXPECT().GetAccount().Times(1).Return(model.Account{
				ActiveCard:     tt.hasCardEnabled,
				AvailableLimit: tt.accountAvailableLimit,
			})

			accountService.EXPECT().AccountEnabled().Times(1).Return(tt.hasAccountEnabled)
			if tt.hasAccountEnabled {
				accountService.EXPECT().AccountCardEnabled().Times(1).Return(tt.hasCardEnabled)
			}

			if tt.hasCardEnabled {
				transactionService.EXPECT().GetTransactionFromOffset(3).Times(1).Return(tt.hasLastThirdTransaction, tt.lastThirdTransaction)
				transactionService.EXPECT().GetTransactionsBy(tt.merchantName, tt.amount, tt.date, 2).Times(1).Return(tt.dupedTransactions)
			}

			if tt.shouldAddTransaction {
				transactionService.EXPECT().AddTransaction(tt.amount, tt.merchantName, tt.date).Times(1)
				accountService.EXPECT().UpdateAvailableLimit(tt.amount).Times(1).Return(tt.expectedAccount)
			}

			transactionUsecase := NewTransaction(transactionService, accountService)
			generatedAccount, generatedViolations := transactionUsecase.CreateTransaction(tt.amount, tt.merchantName, tt.date)
			assert.EqualValues(tt.expectedAccount, generatedAccount, "Account retrieved mismatch")
			assert.EqualValues(tt.expectedViolations, generatedViolations, "Violations mismatch")
		})
	}
}
