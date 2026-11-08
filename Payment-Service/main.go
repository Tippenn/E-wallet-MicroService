package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Tippenn/E-wallet-MicroService/domains/payment/entities"
	"github.com/Tippenn/E-wallet-MicroService/wizards"
)

func main() {
	log.Printf("payment-service started, consuming queue %q", wizards.Config.RabbitMQ.PaymentQueue)

	defer wizards.RabbitMQBroker.Close()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := wizards.PostgresDatabase.GetInstance().AutoMigrate(
		&entities.Payment{},
	)

	if err != nil {
		log.Fatal("migration failed: ", err)
	}
	wizards.StartAllConsumer(ctx)

	log.Println("payment-service started")
	<-ctx.Done()
	log.Println("payment-service shutting down")
}
