package service

import (
	"time"

	"github.com/mickambar19/authorizer/src/helpers"
	"github.com/mickambar19/authorizer/src/model"
)

type TransactionManager struct {
	transactions []model.Transaction
}

func NewTransactionsManager() *TransactionManager {
	return &TransactionManager{}
}

func (t *TransactionManager) AddTransaction(amount int, merchant string, time time.Time) {
	t.transactions = append(t.transactions, model.Transaction{
		Amount:   amount,
		Merchant: merchant,
		Time:     time,
	})
}

func (t *TransactionManager) GetLastTransaction() model.Transaction {
	size := len(t.transactions) - 1
	if size == -1 {
		return model.Transaction{}
	}
	return t.transactions[size]
}

func (t *TransactionManager) GetTransactionFromOffset(endOffset int) (transactionExists bool, transaction model.Transaction) {
	size := len(t.transactions)
	fromIdx := size - endOffset
	if size == 0 || (fromIdx < 0 || fromIdx >= size) {
		return false, model.Transaction{}
	}

	return true, t.transactions[fromIdx]
}

func (t *TransactionManager) GetTransactionsBy(merchant string, amount int, newTime time.Time, withinIntervalMins int) []model.Transaction {
	transactions := []model.Transaction{}

	for i := len(t.transactions) - 1; i >= 0; i-- {
		currentTransaction := t.transactions[i]
		if !helpers.AreDatesWithinInterval(newTime, currentTransaction.Time, withinIntervalMins) {
			break
		}
		if currentTransaction.Merchant == merchant && currentTransaction.Amount == amount {
			transactions = append(transactions, currentTransaction)
		}
	}

	return transactions
}
