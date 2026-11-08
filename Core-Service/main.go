package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	JWTVerification "e-wallet-system-api/domains/JWTVerification/entities"
	transaction "e-wallet-system-api/domains/transactions/entities"
	users "e-wallet-system-api/domains/users/entities"
	wallets "e-wallet-system-api/domains/wallets/entities"
	"e-wallet-system-api/wizards"

	"github.com/gin-gonic/gin"
)

func main() {
	defer wizards.RabbitMQBroker.Close()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err := wizards.PostgresDatabase.GetInstance().AutoMigrate(
		&users.User{},
		&JWTVerification.ServerJWT{},
		&wallets.Wallet{},
		&transaction.Transaction{},
	)

	if err != nil {
		log.Fatal("migration failed: ", err)
	}
	router := gin.Default()

	wizards.RegisterServer(router)
	wizards.StartAllConsumer(ctx)
	router.Run(fmt.Sprintf(":%d", wizards.Config.Server.Port))
}
