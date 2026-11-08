package wizards

import (
	"github.com/Tippenn/E-wallet-MicroService/config"
	serverJWTRepository "github.com/Tippenn/E-wallet-MicroService/domains/JWTVerification/repository"
	serverJWTUsecase "github.com/Tippenn/E-wallet-MicroService/domains/JWTVerification/usecase"
	transactionshttp "github.com/Tippenn/E-wallet-MicroService/domains/transactions/handler/http"
	transactionsRepository "github.com/Tippenn/E-wallet-MicroService/domains/transactions/repository"
	transactionsUsecase "github.com/Tippenn/E-wallet-MicroService/domains/transactions/usecase"
	userhttp "github.com/Tippenn/E-wallet-MicroService/domains/users/handler/http"
	userRepository "github.com/Tippenn/E-wallet-MicroService/domains/users/repository"
	userUsecase "github.com/Tippenn/E-wallet-MicroService/domains/users/usecase"
	walletshttp "github.com/Tippenn/E-wallet-MicroService/domains/wallets/handler/http"
	walletsRepository "github.com/Tippenn/E-wallet-MicroService/domains/wallets/repository"
	walletsUsecase "github.com/Tippenn/E-wallet-MicroService/domains/wallets/usecase"
	"github.com/Tippenn/E-wallet-MicroService/infrastructure/cache"
	"github.com/Tippenn/E-wallet-MicroService/infrastructure/database"
	"github.com/Tippenn/E-wallet-MicroService/infrastructure/messaging"
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
