package response

type PayResponse struct {
	TransactionID string `json:"transactionID" validate:"required"`
	Message       string `json:"message" validate:"required"`
	NewBalance    string `json:"newBalance" validate:"required"`
}
