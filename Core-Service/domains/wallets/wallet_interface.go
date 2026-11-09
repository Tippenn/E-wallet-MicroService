package wallets

import (
	"context"

	"github.com/Tippenn/E-wallet-MicroService/domains/wallets/entities"
	"github.com/Tippenn/E-wallet-MicroService/domains/wallets/models/response"
)

type WalletUseCase interface {
	CheckBalance(ctx context.Context, UserId string) (*response.CheckBalanceResponse, error)
	AddBalanceByUserId(ctx context.Context, UserId string, amount float64) (*entities.Wallet, error)
	DeductBalanceByUserId(ctx context.Context, UserId string, amount float64) (*entities.Wallet, error)
	CreateWallet(ctx context.Context, userId string, balance float64) (*entities.Wallet, error)
}

type WalletRepository interface {
	FindWalletbyUserID(ctx context.Context, Id string) (*entities.Wallet, error)
	CreateWallet(ctx context.Context, Id string, userId string, balance float64) (*entities.Wallet, error)
	UpdateWallet(ctx context.Context, userId string, amount float64) error
	DeleteWallet(ctx context.Context, userId string) error
}
