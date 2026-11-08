package entities

type TransactionType string

const (
	TopUp       TransactionType = "topup"
	TransferIn  TransactionType = "transfer_in"
	TransferOut TransactionType = "transfer_out"
	PaymentIn   TransactionType = "payment_in"
	PaymentOut  TransactionType = "payment_out"
)
