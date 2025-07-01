package infrastructures

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGorm(dsn string) (*gorm.DB, *sql.DB) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("error connecting database : %+v", err.Error()))
	}

	connection, err := db.DB()
	if err != nil {
		panic(err)
	}

	connection.SetMaxIdleConns(5)
	connection.SetConnMaxLifetime(time.Hour)
	connection.SetMaxOpenConns(100)

	return db, connection
}
