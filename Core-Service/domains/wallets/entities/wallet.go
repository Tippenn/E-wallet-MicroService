package entities

type Wallet struct {
	Id      string `gorm:"primaryKey;type:varchar(255)"`
	UserId  string `gorm:"type:varchar(255);index;not null"`
	Balance float64
}
