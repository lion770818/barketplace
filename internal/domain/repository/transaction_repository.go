package repository

import (
	"marketplace/internal/domain/entity"
)

type TransactionRepository interface {
	InsterTransaction(transaction *entity.Transaction) (*entity.Transaction, map[string]string)
	GetTransaction(transactionId string) (*entity.Transaction, error)
	GetTransactionList() ([]entity.Transaction, error)
	UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, map[string]string)
}
