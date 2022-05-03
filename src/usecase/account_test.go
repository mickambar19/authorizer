package usecase

import (
	"testing"

	"github.com/golang/mock/gomock"
	mocks "github.com/mickambar19/authorizer/src/mock"
	"github.com/mickambar19/authorizer/src/model"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -source=account.go -destination=../mock/usecase_account.go -package=mocks

func TestCreateAccount(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name               string
		accountEnabled     bool
		activeCard         bool
		availableLimit     int
		expectedAccount    model.Account
		expectedViolations []model.Violation
	}{
		{
			name:           "Account initialized successfully",
			accountEnabled: false,
			activeCard:     true,
			availableLimit: 1000,
			expectedAccount: model.Account{
				ActiveCard:     true,
				AvailableLimit: 1000,
			},
			expectedViolations: []model.Violation{},
		},
		{
			name:           "Account already initialized",
			accountEnabled: true,
			activeCard:     false,
			availableLimit: 200,
			expectedAccount: model.Account{
				ActiveCard:     true,
				AvailableLimit: 1000,
			},
			expectedViolations: []model.Violation{"account-already-initialized"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			accountManager := mocks.NewMockAccountManager(ctrl)

			accountManager.EXPECT().AccountEnabled().Return(tt.accountEnabled)

			if tt.accountEnabled {
				accountManager.EXPECT().GetAccount().Times(1).Return(tt.expectedAccount)
			} else {
				accountManager.EXPECT().GetAccount().Times(0)
				accountManager.EXPECT().CreateAccount(tt.activeCard, tt.availableLimit).Times(1).Return(model.Account{ActiveCard: tt.activeCard, AvailableLimit: tt.availableLimit})
			}
			accountUsecase := NewAccount(accountManager)
			generatedAccount, generatedViolations := accountUsecase.CreateAccount(tt.activeCard, tt.availableLimit)
			assert.EqualValues(tt.expectedAccount, generatedAccount, "Account retrieved mismatch")
			assert.EqualValues(tt.expectedViolations, generatedViolations, "Violations mismatch")
		})
	}
}
