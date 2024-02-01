package repositories

import (
	"errors"
	"marketplace/internal/domain/entity"
	"marketplace/internal/domain/repository"
	"strings"

	"github.com/jinzhu/gorm"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{db}
}

// TransactionRepo implements the repository.TransactionRepository interface
var _ repository.TransactionRepository = &TransactionRepo{}

func (r *TransactionRepo) InsterTransaction(transaction *entity.Transaction) (*entity.Transaction, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&transaction).Error
	if err != nil {
		//If the transactionId is already taken
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["transactionId_taken"] = "transactionId already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return transaction, nil
}

func (t *TransactionRepo) GetTransaction(transactionId string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := t.db.Debug().Where("transaction_id = ?", transactionId).Take(&transaction).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("transaction not found")
	}
	return &transaction, nil
}

func (r *TransactionRepo) GetTransactionList() ([]entity.Transaction, error) {
	var transactionList []entity.Transaction
	err := r.db.Debug().Find(&transactionList).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("transaction not found")
	}
	return transactionList, nil
}

func (r *TransactionRepo) UpdateTransaction(transaction *entity.Transaction) (*entity.Transaction, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&transaction).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return transaction, nil
}
