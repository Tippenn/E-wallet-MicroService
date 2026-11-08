package entities

import (
	"time"
)

type ServerJWT struct {
	Id        string    `gorm:"primaryKey;type:varchar(255)"`
	UserId    string    `gorm:"type:varchar(255);index;not null"`
	Revoked   bool      `gorm:"default:false;not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
}
