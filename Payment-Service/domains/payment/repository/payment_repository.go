package repository

import (
	"context"
	"time"

	"github.com/Tippenn/E-wallet-MicroService/domains/payment"
	"github.com/Tippenn/E-wallet-MicroService/domains/payment/entities"
	"github.com/Tippenn/E-wallet-MicroService/infrastructure/database"
)

type paymentRepository struct {
	db database.Database
}

func NewPaymentRepository(db database.Database) payment.PaymentRepository {
	return &paymentRepository{db: db}
}

// CreatePayment implements [payment.PaymentRepository].
func (p *paymentRepository) CreatePayment(ctx context.Context, paymentID string, userID string, amount float64) error {
	payment := &entities.Payment{
		PaymentId:   paymentID,
		UserId:      userID,
		Amount:      amount,
		ProcessedAt: time.Now().UTC(),
	}

	result := p.db.GetInstance().Create(payment)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
