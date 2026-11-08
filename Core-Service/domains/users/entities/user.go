package entities

import (
	ServerJWT "github.com/Tippenn/E-wallet-MicroService/domains/JWTVerification/entities"
	Wallet "github.com/Tippenn/E-wallet-MicroService/domains/wallets/entities"
)

type User struct {
	Id         string `gorm:"primaryKey;type:varchar(255)"`
	Email      string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Name       string
	Password   string
	Enabled    int
	ServerJWTs []ServerJWT.ServerJWT `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Wallet     Wallet.Wallet         `gorm:"foreignKey:UserId;"`
}
