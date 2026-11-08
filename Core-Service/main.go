package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	JWTVerification "github.com/Tippenn/E-wallet-MicroService/domains/JWTVerification/entities"
	transaction "github.com/Tippenn/E-wallet-MicroService/domains/transactions/entities"
	users "github.com/Tippenn/E-wallet-MicroService/domains/users/entities"
	wallets "github.com/Tippenn/E-wallet-MicroService/domains/wallets/entities"
	"github.com/Tippenn/E-wallet-MicroService/wizards"

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
