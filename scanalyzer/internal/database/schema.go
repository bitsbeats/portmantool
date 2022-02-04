// Copyright 2020-2022 Thomann Bits & Beats GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	ScanID	time.Time	`gorm:"primaryKey;index"`
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
