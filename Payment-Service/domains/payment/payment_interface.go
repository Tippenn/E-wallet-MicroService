package payment

import "context"

type PaymentUseCase interface {
	ProcessPayment(ctx context.Context, body []byte) error
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, paymentID string, userID string, amount float64) error
}
