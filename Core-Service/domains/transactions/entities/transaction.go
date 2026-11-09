package entities

import (
	"time"
)

type Transaction struct {
	Id             string            `gorm:"primaryKey" json:"id"`
	IdemptonecyKey *string            `gorm:"uniqueIndex"`
	Type           TransactionType   `json:"type"`
	Amount         float64           `json:"amount"`
	From           *string           `json:"from"`
	To             string            `json:"to"`
	Timestamp      time.Time         `json:"timestamp"`
	Status         TransactionStatus `json:"status"`
}
