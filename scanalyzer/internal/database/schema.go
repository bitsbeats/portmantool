package database

import "time"

type Target struct {
	Address		string	`gorm:"type:inet;primaryKey"`
	Port		uint16	`gorm:"check:port > 0 AND port < 65536;primaryKey;autoIncrement:false"`
	Protocol	string	`gorm:"type:protocol;primaryKey"`
}

type ActualState struct {
	Target
	State	string		`gorm:"type:state;not null"`
	ScanID	time.Time	`gorm:"primaryKey"`
}

type ExpectedState struct {
	Target
	State	string	`gorm:"type:state;not null"`
	Comment	string	`gorm:"not null"`
}

type Scan struct {
	ID	time.Time
	Ports	[]ActualState	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
