package wizards

import (
	"e-wallet-system-api/config"
	paymentRepository "e-wallet-system-api/domains/payment/repository"
	paymentUsecase "e-wallet-system-api/domains/payment/usecase"
	"e-wallet-system-api/infrastructure/database"
	"e-wallet-system-api/infrastructure/messaging"
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
