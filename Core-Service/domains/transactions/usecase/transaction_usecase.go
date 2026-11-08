package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Tippenn/E-wallet-MicroService/domains/transactions"
	"github.com/Tippenn/E-wallet-MicroService/domains/transactions/entities"
	"github.com/Tippenn/E-wallet-MicroService/domains/transactions/models/response"
	"github.com/Tippenn/E-wallet-MicroService/domains/users"
	"github.com/Tippenn/E-wallet-MicroService/domains/wallets"
	"github.com/Tippenn/E-wallet-MicroService/infrastructure/messaging"
	"github.com/Tippenn/E-wallet-MicroService/shared/dto/message"

	"github.com/google/uuid"
)

type transactionUseCase struct {
	repo              transactions.TransactionRepository
	transactor        transactions.Transactor
	messageBroker     messaging.MessageBroker
	walletUc          wallets.WalletUseCase
	userUc            users.UserUseCase
	paymentRoutingKey string
}

func NewTransactionUseCase(repo transactions.TransactionRepository, transactor transactions.Transactor, walletUc wallets.WalletUseCase, userUc users.UserUseCase, messageBroker messaging.MessageBroker, paymentRoutingKey string) transactions.TransactionUseCase {
	return transactionUseCase{repo: repo, transactor: transactor, walletUc: walletUc, userUc: userUc, messageBroker: messageBroker, paymentRoutingKey: paymentRoutingKey}
}

// CompletePaymentSession implements [transactions.TransactionUseCase].
func (t transactionUseCase) CompletePaymentSession(ctx context.Context, userId string, merchantId string, amount float64, idemptKey string) (*response.PayResponse, error) {

	transactionID := uuid.New().String()
	//create payment out and deduct
	t.transactor.WithTx(ctx, func(transactionRepo transactions.TransactionRepository, walletRepo wallets.WalletRepository) error {
		//payment out created
		err := transactionRepo.CreateTransaction(ctx,
			transactionID,
			&idemptKey,
			amount,
			&userId,
			merchantId,
			entities.Pending,
			entities.PaymentOut,
		)
		if err != nil {
			return err
		}
		//wallet deducted
		err = walletRepo.UpdateWallet(ctx, userId, -amount)
		if err != nil {
			return err
		}
		return nil
	})
	//publish message
	body, err := json.Marshal(message.PaymentMessage{
		PaymentID: transactionID,
		UserID:    userId,
		Amount:    amount,
		CreatedAt: time.Now().UTC(),
	})

	if err != nil {
		t.walletUc.AddBalanceByUserId(ctx, userId, amount)
		return nil, err
	}
	err = t.messageBroker.Publish(ctx, t.paymentRoutingKey, body)
	if err != nil {
		t.walletUc.AddBalanceByUserId(ctx, userId, amount)
		_ = t.repo.UpdateTransactionStatus(ctx, transactionID, entities.Failed)
		return nil, errors.New("could not queue payment, please try again")
	}

	wallet, err := t.walletUc.CheckBalance(ctx, userId)
	return &response.PayResponse{
		TransactionID: transactionID,
		Message:       "Payment Pending",
		NewBalance:    wallet.Balance,
	}, nil
}

// HandlePaymentStatus implements [transactions.TransactionUseCase].
func (t transactionUseCase) HandlePaymentStatus(ctx context.Context, body []byte) error {
	var message *message.PaymentStatusMessage
	err := json.Unmarshal(body, &message)
	if err != nil {
		return err
	}

	//translate the queue's vocabulary into your own status
	var newStatus entities.TransactionStatus
	switch message.Status {
	case "success":
		newStatus = entities.Completed
	case "failed":
		newStatus = entities.Failed
	case "expired":
		newStatus = entities.Expired
	default:
		return fmt.Errorf("unknown status %q", message.Status)
	}

	// load our own record (we need the user and amount for a refund)
	transaction, err := t.repo.FindTransactionbyId(ctx, message.PaymentID)
	if err != nil {
		return err
	}

	err = t.repo.UpdateTransactionStatus(ctx, message.PaymentID, newStatus)
	if err != nil {
		return err
	}
	if newStatus != entities.Completed {
		t.walletUc.AddBalanceByUserId(ctx, *transaction.From, transaction.Amount)
	}
	return nil
}

// ShowHistory implements [transactions.TransactionUseCase].
func (t transactionUseCase) ShowHistory(ctx context.Context, userId string, limit int, offset int) ([]*entities.Transaction, int, error) {
	allTransaction, err := t.repo.FindAllTransactionbyUserId(ctx, userId)
	if err != nil {
		return nil, 0, err
	}

	total := len(allTransaction)

	// bounds check so a page past the end returns [] instead of panicking
	if offset >= total {
		return []*entities.Transaction{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}

	return allTransaction[offset:end], total, nil
}

// TopUp implements [transactions.TransactionUseCase].
func (t transactionUseCase) TopUp(ctx context.Context, userID string, amount float64) (*response.TopUpResponse, error) {
	wallet, err := t.walletUc.AddBalanceByUserId(ctx, userID, amount)
	if err != nil {
		return nil, err

	}
	topupId := uuid.New().String()
	t.repo.CreateTransaction(ctx,
		topupId,
		nil,
		amount,
		nil,
		userID,
		entities.Completed,
		entities.TopUp,
	)
	return &response.TopUpResponse{
		TransactionID: topupId,
		Message:       "Top-up Successful",
		NewBalance:    fmt.Sprintf("%.1f", wallet.Balance),
	}, nil
}

// Transfer implements [transactions.TransactionUseCase].
func (t transactionUseCase) Transfer(ctx context.Context, fromID string, toID string, amount float64) (*response.TransferResponse, error) {
	_, err := t.walletUc.DeductBalanceByUserId(ctx, fromID, amount)
	if err != nil {
		return nil, err
	}
	transferId := uuid.New().String()

	_, err = t.walletUc.AddBalanceByUserId(ctx, toID, amount)
	if err != nil {
		return nil, err
	}
	t.repo.CreateTransaction(ctx,
		transferId,
		nil,
		amount,
		&fromID,
		toID,
		entities.Completed,
		entities.TransferIn,
	)
	return &response.TransferResponse{
		TransactionID: transferId,
		Status:        entities.TransactionStatus(entities.Pending),
		Message:       "Transfer Initiated Successfully",
	}, nil
}
