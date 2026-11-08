package message

import "time"

type PaymentMessage struct {
	PaymentID string    `json:"payment_id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
