package http

import (
	"github.com/Tippenn/E-wallet-MicroService/domains/wallets"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Userhttp struct {
	uc wallets.WalletUseCase
}

func NewWalletHttp(uc wallets.WalletUseCase) *Userhttp {
	return &Userhttp{uc: uc}
}

func (handler *Userhttp) CheckBalance(c *gin.Context) {
	ctx := c.Request.Context()
	userId := c.GetString("user_Id")

	validate := validator.New()
	err := validate.StructCtx(ctx, userId)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	checkBalanceResponse, err := handler.uc.CheckBalance(ctx, userId)

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"userID":  checkBalanceResponse.UserId,
		"balance": checkBalanceResponse.Balance,
	})
}
