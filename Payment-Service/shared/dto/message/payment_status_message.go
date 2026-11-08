package message

import (
	"time"
)

type PaymentStatusMessage struct {
	PaymentID   string    `json:"payment_id"`
	Status      string    `json:"status"` // success | failed | expired
	Reason      string    `json:"reason,omitempty"`
	ProcessedAt time.Time `json:"processed_at"`
}
