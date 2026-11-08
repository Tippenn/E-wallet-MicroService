package wizards

import (
	"e-wallet-system-api/config"
	serverJWTRepository "e-wallet-system-api/domains/JWTVerification/repository"
	serverJWTUsecase "e-wallet-system-api/domains/JWTVerification/usecase"
	transactionshttp "e-wallet-system-api/domains/transactions/handler/http"
	transactionsRepository "e-wallet-system-api/domains/transactions/repository"
	transactionsUsecase "e-wallet-system-api/domains/transactions/usecase"
	userhttp "e-wallet-system-api/domains/users/handler/http"
	userRepository "e-wallet-system-api/domains/users/repository"
	userUsecase "e-wallet-system-api/domains/users/usecase"
	walletshttp "e-wallet-system-api/domains/wallets/handler/http"
	walletsRepository "e-wallet-system-api/domains/wallets/repository"
	walletsUsecase "e-wallet-system-api/domains/wallets/usecase"
	"e-wallet-system-api/infrastructure/cache"
	"e-wallet-system-api/infrastructure/database"
	"e-wallet-system-api/infrastructure/messaging"
)

var (
	Config = config.GetConfig()

	//Storage
	PostgresDatabase = database.NewPostgresDatabase(Config)
	RedisCache       = cache.NewRedisDatabase(Config)

	//message broker
	RabbitMQBroker = messaging.NewRabbitMQBroker(Config)

	//transactor
	TransactionTransactor = transactionsRepository.NewTransactor(PostgresDatabase)

	//repo
	UserRepo        = userRepository.NewUserRepository(PostgresDatabase)
	ServerJWTRepo   = serverJWTRepository.NewServerJWTRepository(RedisCache)
	WalletRepo      = walletsRepository.NewWalletRepository(PostgresDatabase)
	TransactionRepo = transactionsRepository.NewTransactionRepository(PostgresDatabase)

	//usecase
	ServerJWTUc   = serverJWTUsecase.NewServerJWTUseCase(ServerJWTRepo)
	UserUc        = userUsecase.NewUserUseCase(UserRepo, ServerJWTUc)
	WalletUc      = walletsUsecase.NewWalletUseCase(WalletRepo)
	TransactionUc = transactionsUsecase.NewTransactionUseCase(TransactionRepo, TransactionTransactor, WalletUc, UserUc, RabbitMQBroker, Config.RabbitMQ.PaymentRoutingKey)

	//http
	UserHttp        = userhttp.NewUserHttp(UserUc)
	TransactionHttp = transactionshttp.NewTransactionHttp(TransactionUc)
	WalletHttp      = walletshttp.NewWalletHttp(WalletUc)
)
