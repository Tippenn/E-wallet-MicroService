package wizards

import (
	"context"
	"log"
)

func StartAllConsumer(ctx context.Context) {
	go func() {
		err := RabbitMQBroker.Consume(ctx, Config.RabbitMQ.StatusQueue, TransactionUc.HandlePaymentStatus)
		if err != nil {
			log.Println("status consumer stopped:", err)
		}
	}()
}
