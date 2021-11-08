package database

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

func Connect(host, user, password, name string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", host, user, password, name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Scan{}, &ActualState{}, &ExpectedState{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
