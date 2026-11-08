package database

import "gorm.io/gorm"

type Database interface {
	GetInstance() *gorm.DB
}
