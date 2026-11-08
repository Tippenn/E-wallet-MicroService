package response

import "e-wallet-system-api/domains/transactions/entities"

type TransactionHistoryResponse struct {
	Transactions []entities.Transaction `json:"transactions"`
	Pagination   entities.Pagination    `json:"pagination"`
}
