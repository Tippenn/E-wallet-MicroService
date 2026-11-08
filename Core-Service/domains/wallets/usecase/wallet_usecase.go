package usecase

import (
	"context"
	"e-wallet-system-api/domains/wallets"
	"e-wallet-system-api/domains/wallets/entities"
	"e-wallet-system-api/domains/wallets/models/response"
	"errors"
	"fmt"
)

type walletUseCase struct {
	repo wallets.WalletRepository
}

func NewWalletUseCase(repo wallets.WalletRepository) wallets.WalletUseCase {
	return walletUseCase{repo: repo}
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

	wallet, err := w.repo.FindWalletbyUserID(ctx, Id)
	if err != nil {
		return nil, err
	}

	w.repo.UpdateWallet(ctx, Id, amount)

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

	w.repo.UpdateWallet(ctx, Id, -amount)

	return wallet, nil
}
