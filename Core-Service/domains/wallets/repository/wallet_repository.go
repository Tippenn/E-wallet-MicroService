package repository

import (
	"context"
	"e-wallet-system-api/domains/wallets"
	"e-wallet-system-api/domains/wallets/entities"
	"e-wallet-system-api/infrastructure/database"
	"errors"

	"gorm.io/gorm"
)

type walletRepository struct {
	db database.Database
}

func NewWalletRepository(db database.Database) wallets.WalletRepository {
	return walletRepository{db: db}
}

// CreateWallet implements [wallets.WalletRepository].
func (w walletRepository) CreateWallet(ctx context.Context, Id string, userId string, balance float64) (*entities.Wallet, error) {
	wallet := &entities.Wallet{
		Id:      Id,
		UserId:  userId,
		Balance: balance,
	}

	result := w.db.GetInstance().Create(wallet)
	if result.Error != nil {
		return nil, result.Error
	}

	return wallet, nil
}

// DeleteWallet implements [wallets.WalletRepository].
func (w walletRepository) DeleteWallet(ctx context.Context, userId string) error {
	panic("unimplemented")
}

// FindWalletbyUserID implements [wallets.WalletRepository].
func (w walletRepository) FindWalletbyUserID(ctx context.Context, UserId string) (*entities.Wallet, error) {
	var wallet *entities.Wallet
	result := w.db.GetInstance().First(&wallet, "user_id=?", UserId)

	if result.Error != nil {
		return nil, result.Error
	}

	return wallet, nil
}

// UpdateWallet implements [wallets.WalletRepository].
func (w walletRepository) UpdateWallet(ctx context.Context, userId string, amount float64) error {
	result := w.db.GetInstance().Model(&entities.Wallet{}).
		Where("user_id = ?", userId).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("wallet not found")
	}

	return nil
}
