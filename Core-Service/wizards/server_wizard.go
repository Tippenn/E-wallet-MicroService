package wizards

import (
	"github.com/Tippenn/E-wallet-MicroService/shared/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.Engine) {

	api := router.Group("/api/v1")
	{
		//User
		api.POST("/users/login", UserHttp.Login)
		api.POST("/users/register", UserHttp.Register)
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(ServerJWTUc))

	//Wallet
	protected.GET("/wallet/balance", WalletHttp.CheckBalance)
	//transaction
	protected.POST("/transactions/topup", TransactionHttp.TopUp)
	protected.POST("/transactions/transfer", TransactionHttp.Transfer)
	protected.POST("/transactions/pay", TransactionHttp.Pay)
	protected.POST("/transactions/history", TransactionHttp.History)

}
