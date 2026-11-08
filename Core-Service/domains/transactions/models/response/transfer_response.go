package response

import (
	"github.com/Tippenn/E-wallet-MicroService/domains/transactions/entities"
)

type TransferResponse struct {
	TransactionID string                     `json:"transactionID" validate:"required"`
	Status        entities.TransactionStatus `json:"status" validate:"required"`
	Message       string                     `json:"message" validate:"required"`
}
