package transactions

import (
	"context"

	"github.com/Tippenn/E-wallet-MicroService/domains/transactions/entities"
	"github.com/Tippenn/E-wallet-MicroService/domains/transactions/models/response"
	"github.com/Tippenn/E-wallet-MicroService/domains/wallets"
)

type TransactionUseCase interface {
	ShowHistory(ctx context.Context, userId string, limit int, offset int) ([]*entities.Transaction, int, error)
	TopUp(ctx context.Context, userID string, amount float64) (*response.TopUpResponse, error)
	Transfer(ctx context.Context, fromID string, toUserID string, amount float64) (*response.TransferResponse, error)
	CompletePaymentSession(ctx context.Context, userID string, merchantID string, amount float64, idemptKey string) (*response.PayResponse, error)
	HandlePaymentStatus(ctx context.Context, body []byte) error
}

type TransactionRepository interface {
	FindTransactionbyId(ctx context.Context, transactionID string) (*entities.Transaction, error)
	FindAllTransactionbyUserId(ctx context.Context, userID string) ([]*entities.Transaction, error)
	CreateTransaction(ctx context.Context, transactionID string, idemptKey *string, amount float64, from *string, to string, status entities.TransactionStatus, types entities.TransactionType) error
	UpdateTransactionStatus(ctx context.Context, transactionID string, status entities.TransactionStatus) error
}

type Transactor interface {
	WithTx(ctx context.Context, fn func(transactionRepo TransactionRepository, walletRepo wallets.WalletRepository) error) error
}
