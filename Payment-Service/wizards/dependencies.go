package wizards

import (
	"github.com/Tippenn/E-wallet-MicroService/config"
	paymentRepository "github.com/Tippenn/E-wallet-MicroService/domains/payment/repository"
	paymentUsecase "github.com/Tippenn/E-wallet-MicroService/domains/payment/usecase"
	"github.com/Tippenn/E-wallet-MicroService/infrastructure/database"
	"github.com/Tippenn/E-wallet-MicroService/infrastructure/messaging"
)

var (
	Config = config.GetConfig()

	//Storage
	PostgresDatabase = database.NewPostgresDatabase(Config)

	//message broker
	RabbitMQBroker = messaging.NewRabbitMQBroker(Config)

	//repo
	PaymentRepo = paymentRepository.NewPaymentRepository(PostgresDatabase)
	//usecase
	PaymentUc = paymentUsecase.NewPaymentUsecase(PaymentRepo, RabbitMQBroker, Config.RabbitMQ.StatusRoutingKey)
)
