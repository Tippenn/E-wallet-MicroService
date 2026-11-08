package repository

import (
	"context"
	"e-wallet-system-api/domains/transactions"
	"e-wallet-system-api/domains/wallets"
	walletRepo "e-wallet-system-api/domains/wallets/repository"
	"e-wallet-system-api/infrastructure/database"

	"gorm.io/gorm"
)

type transactionTransactor struct {
	db database.Database
}

type transactorDB struct {
	db *gorm.DB
}

func (d transactorDB) GetInstance() *gorm.DB {
	return d.db
}
func NewTransactor(db database.Database) transactions.Transactor {
	return &transactionTransactor{db: db}
}

func (t *transactionTransactor) WithTx(ctx context.Context, fn func(transactionRepo transactions.TransactionRepository, walletRepo wallets.WalletRepository) error) error {
	return t.db.GetInstance().Transaction(func(tx *gorm.DB) error {
		txdb := &transactorDB{db: tx}
		return fn(NewTransactionRepository(txdb), walletRepo.NewWalletRepository(txdb))
	})
}
