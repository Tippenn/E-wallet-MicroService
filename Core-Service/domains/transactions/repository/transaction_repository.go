package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Tippenn/E-wallet-MicroService/domains/transactions"
	"github.com/Tippenn/E-wallet-MicroService/domains/transactions/entities"
	"github.com/Tippenn/E-wallet-MicroService/infrastructure/database"
)

type transactionRepository struct {
	db database.Database
}

func NewTransactionRepository(db database.Database) transactions.TransactionRepository {
	return transactionRepository{db: db}
}

// CreateTransaction implements [transactions.TransactionRepository].
func (t transactionRepository) CreateTransaction(ctx context.Context, transactionId string, idemptKey *string, amount float64, from *string, to string, status entities.TransactionStatus, types entities.TransactionType) error {
	transaction := &entities.Transaction{
		Id:             transactionId,
		IdemptonecyKey: *idemptKey,
		Amount:         amount,
		From:           from,
		To:             to,
		Timestamp:      time.Now().UTC(),
		Status:         status,
		Type:           types,
	}

	result := t.db.GetInstance().Create(transaction)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindTransactionbyId implements [transactions.TransactionRepository].
func (t transactionRepository) FindTransactionbyId(ctx context.Context, transactionId string) (*entities.Transaction, error) {
	var transaction *entities.Transaction
	result := t.db.GetInstance().First(&transaction, "id=?", transactionId)

	if result.Error != nil {

		return nil, result.Error
	}

	return transaction, nil
}

// FindAllTransactionbyUserId implements [transactions.TransactionRepository].
func (t transactionRepository) FindAllTransactionbyUserId(ctx context.Context, userId string) ([]*entities.Transaction, error) {
	var transactions []*entities.Transaction
	err := t.db.GetInstance().
		Where("from_id = ? OR to_id = ?", userId, userId).
		Order("timestamp DESC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// UpdateTransactionStatus implements [transactions.TransactionRepository].
func (t transactionRepository) UpdateTransactionStatus(ctx context.Context, transactionId string, status entities.TransactionStatus) error {
	result := t.db.GetInstance().Model(&entities.Transaction{}).
		Where("id = ?", transactionId).
		Updates(map[string]interface{}{
			"status": entities.Completed,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("wallet not found")
	}

	return nil
}
