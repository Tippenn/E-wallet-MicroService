package http

import (
	"net/http"
	"strconv"

	"github.com/Tippenn/E-wallet-MicroService/domains/transactions"
	"github.com/Tippenn/E-wallet-MicroService/domains/transactions/models/request"
	"github.com/Tippenn/E-wallet-MicroService/shared/dto/responses"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TransactionHttp struct {
	uc transactions.TransactionUseCase
}

func NewTransactionHttp(uc transactions.TransactionUseCase) *TransactionHttp {
	return &TransactionHttp{uc: uc}
}

func (handler *TransactionHttp) TopUp(c *gin.Context) {
	ctx := c.Request.Context()
	requestBody := &request.TopUpRequest{}
	userId := c.GetString("user_Id")

	err := c.ShouldBindJSON(&requestBody)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Invalid request",
		})
		return
	}

	validate := validator.New()
	err = validate.StructCtx(ctx, requestBody)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	topupResponse, err := handler.uc.TopUp(ctx, userId, requestBody.Amount)

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"transactionID": topupResponse.TransactionID,
		"message":       "TopUp successful!",
		"newBalance":    topupResponse.NewBalance,
	})
}

func (handler *TransactionHttp) Transfer(c *gin.Context) {
	ctx := c.Request.Context()
	requestBody := &request.TransferRequest{}
	userId := c.GetString("user_Id")

	err := c.ShouldBindJSON(&requestBody)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Invalid request",
		})
		return
	}

	validate := validator.New()
	err = validate.StructCtx(ctx, requestBody)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	transferResponse, err := handler.uc.Transfer(ctx, userId, requestBody.ToUserID, requestBody.Amount)

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"transactionID": transferResponse.TransactionID,
		"status":        transferResponse.Status,
		"message":       "TopUp successful!",
	})
}

func (handler *TransactionHttp) Pay(c *gin.Context) {
	ctx := c.Request.Context()
	requestBody := &request.PayRequest{}
	userId := c.GetString("user_Id")
	idempotencyKey := c.GetHeader("Idempotency-Key")
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header is required"})
		return
	}
	err := c.ShouldBindJSON(&requestBody)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Invalid request",
		})
		return
	}

	validate := validator.New()
	err = validate.StructCtx(ctx, requestBody)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	payResponse, err := handler.uc.CompletePaymentSession(ctx, userId, requestBody.MerchantID, requestBody.Amount, idempotencyKey)

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"transactionID": payResponse.TransactionID,
		"message":       "TopUp successful!",
		"newBalance":    payResponse.NewBalance,
	})
}

func (handler *TransactionHttp) History(c *gin.Context) {
	ctx := c.Request.Context()
	userId := c.GetString("user_Id")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, responses.BasicResponses{Error: "unauthorized"})
		return
	}

	// Defaults from the spec: page=1, limit=10
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit
	transactions, total, err := handler.uc.ShowHistory(ctx, userId, limit, offset)

	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"transactions": transactions,
		"pagination": gin.H{
			"currentPage": page,
			"totalPages":  (total + limit - 1) / limit,
			"totalItems":  total,
		}})
}
