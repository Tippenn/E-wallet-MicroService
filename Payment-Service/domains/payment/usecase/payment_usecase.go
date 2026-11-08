package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Tippenn/E-wallet-MicroService/domains/payment"
	"github.com/Tippenn/E-wallet-MicroService/infrastructure/messaging"
	"github.com/Tippenn/E-wallet-MicroService/shared/dto/message"
)

type paymentUseCase struct {
	repo             payment.PaymentRepository
	messageBroker    messaging.MessageBroker
	statusRoutingKey string
}

func NewPaymentUsecase(repo payment.PaymentRepository, messageBroker messaging.MessageBroker, statusRoutingKey string) payment.PaymentUseCase {
	return &paymentUseCase{repo: repo, messageBroker: messageBroker, statusRoutingKey: statusRoutingKey}
}

// ProcessPayment implements [payment.PaymentUseCase].
func (p *paymentUseCase) ProcessPayment(ctx context.Context, body []byte) error {
	var msg *message.PaymentMessage
	const paymentTTL = 10 * time.Minute
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return err
	}

	//check if expired or not
	if time.Since(msg.CreatedAt) > paymentTTL {
		body, err = json.Marshal(message.PaymentStatusMessage{
			PaymentID: msg.PaymentID,
			Status:    "expired",
			Reason:    "Payment failed due to expired session",
		})

		err = p.messageBroker.Publish(ctx, p.statusRoutingKey, body)
		if err != nil {
			return err
		}
		return nil
	}

	//simulate payment processing
	time.Sleep(time.Second * 3)

	//create payment log inside database repository
	err = p.repo.CreatePayment(ctx, msg.PaymentID, msg.UserID, msg.Amount)
	if err != nil {
		return err
	}

	//payment success
	body, err = json.Marshal(message.PaymentStatusMessage{
		PaymentID: msg.PaymentID,
		Status:    "success",
		Reason:    "Payment Completed Successfully",
	})

	err = p.messageBroker.Publish(ctx, p.statusRoutingKey, body)
	if err != nil {
		return err
	}

	return nil
}
