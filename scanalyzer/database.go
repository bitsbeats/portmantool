package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

type Target struct {
	Address		string	`gorm:"type:inet;primaryKey"`
	Port		uint16	`gorm:"check:port > 0 AND port < 65536;primaryKey;autoIncrement:false"`
	Protocol	string	`gorm:"type:protocol;primaryKey"`
}

type ActualState struct {
	Target
	State	string	`gorm:"type:state;not null"`
	ScanID	int32	`gorm:"primaryKey;autoIncrement:false"`
	Scan	Scan	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ExpectedState struct {
	Target
	State	string	`gorm:"type:state;not null"`
	Comment	string	`gorm:"not null"`
}

type Scan struct {
	ID		int32
	Timestamp	time.Time	`gorm:"not null"`
}

func InitDatabase() *gorm.DB {
	password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		log.Fatal("DB_PASSWORD is not set")
	}

	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		host = "localhost"
	}

	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		user = "postgres"
	}

	name, ok := os.LookupEnv("DB_NAME")
	if !ok {
		name = "postgres"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", host, user, password, name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
