package entities

import "time"

type Payment struct {
	PaymentId   string
	UserId      string
	Amount      float64
	ProcessedAt time.Time
}
