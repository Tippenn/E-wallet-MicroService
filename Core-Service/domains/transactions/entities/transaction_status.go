package entities

type TransactionStatus string

const (
	Completed TransactionStatus = "complete"
	Pending   TransactionStatus = "pending"
	Failed    TransactionStatus = "failed"
	Expired   TransactionStatus = "expired"
)
