package entities

type TransactionType string

const (
	TopUp      TransactionType = "topup"
	Transfer   TransactionType = "transfer"
	PaymentIn  TransactionType = "payment_in"
	PaymentOut TransactionType = "payment_out"
)
