package entities

type PaymentStatus string

const (
	Completed PaymentStatus = "complete"
	Pending   PaymentStatus = "pending"
	Failed    PaymentStatus = "failed"
	Expired   PaymentStatus = "expired"
)
