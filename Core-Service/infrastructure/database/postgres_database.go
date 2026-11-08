package database

import (
	"fmt"
	"sync"

	"e-wallet-system-api/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase struct {
	db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *PostgresDatabase
)

func NewPostgresDatabase(conf *config.Config) Database {
	once.Do(func() {
		db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d sslmode=%s timezone=%s",
			conf.Database.Host,
			conf.Database.DbName,
			conf.Database.User,
			conf.Database.Password,
			conf.Database.Port,
			conf.Database.SslMode,
			conf.Database.TimeZone)),
			&gorm.Config{},
		)

		if err != nil {
			panic(err)
		}

		dbInstance = &PostgresDatabase{
			db: db,
		}
	})
	return dbInstance
}

// GetInstance implements [Database].
func (p *PostgresDatabase) GetInstance() *gorm.DB {
	return dbInstance.db
}
