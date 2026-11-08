package response

import "github.com/Tippenn/E-wallet-MicroService/domains/transactions/entities"

type TransactionHistoryResponse struct {
	Transactions []entities.Transaction `json:"transactions"`
	Pagination   entities.Pagination    `json:"pagination"`
}
