package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tippenn/E-wallet-MicroService/domains/wallets"
	"github.com/Tippenn/E-wallet-MicroService/domains/wallets/entities"
	"github.com/Tippenn/E-wallet-MicroService/domains/wallets/models/response"
	"github.com/google/uuid"
)

type walletUseCase struct {
	repo wallets.WalletRepository
}

func NewWalletUseCase(repo wallets.WalletRepository) wallets.WalletUseCase {
	return walletUseCase{repo: repo}
}

// CreateWallet implements [wallets.WalletUseCase].
func (w walletUseCase) CreateWallet(ctx context.Context, userId string, balance float64) (*entities.Wallet, error) {
	wid := uuid.New().String()
	wallet, err := w.repo.CreateWallet(ctx, wid, userId, balance)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

// CheckBalance implements [wallets.WalletUseCase].
func (w walletUseCase) CheckBalance(ctx context.Context, Id string) (*response.CheckBalanceResponse, error) {
	wallet, err := w.repo.FindWalletbyUserID(ctx, Id)
	if err != nil {
		return nil, err
	}

	return &response.CheckBalanceResponse{
		UserId:  wallet.UserId,
		Balance: fmt.Sprintf("%.1f", wallet.Balance),
	}, nil
}

// AddBalance implements [wallets.WalletUseCase].
func (w walletUseCase) AddBalanceByUserId(ctx context.Context, Id string, amount float64) (*entities.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	_, err := w.repo.FindWalletbyUserID(ctx, Id)
	if err != nil {
		return nil, err
	}

	err = w.repo.UpdateWallet(ctx, Id, amount)
	if err != nil {
		return nil, err
	}

	wallet, err := w.repo.FindWalletbyUserID(ctx, Id)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// DeductBalance implements [wallets.WalletUseCase].
func (w walletUseCase) DeductBalanceByUserId(ctx context.Context, Id string, amount float64) (*entities.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	wallet, err := w.repo.FindWalletbyUserID(ctx, Id)
	if err != nil {
		return nil, err
	}

	//cek apakah ada uang
	if wallet.Balance <= amount {
		return nil, errors.New("insufficient balance")
	}

	err = w.repo.UpdateWallet(ctx, Id, -amount)
	if err != nil {
		return nil, err
	}

	wallet, err = w.repo.FindWalletbyUserID(ctx, Id)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}
