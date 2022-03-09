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

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	FailedImports = "failed_imports"
	LastImport = "last_successful_import"
)

func Retrieve(db *gorm.DB, keys ...string) (info map[string]string, err error) {
	result := make([]Info, 0)

	err = db.Where("key IN ?", keys).Find(&result).Error
	if err != nil {
		return nil, err
	}

	info = make(map[string]string)

	for _, i := range result {
		info[i.Key] = i.Value
	}

	return info, nil
}

func Persist(db *gorm.DB, info map[string]string) error {
	pairs := make([]Info, 0)

	for k, v := range info {
		pairs = append(pairs, Info{
			Key: k,
			Value: v,
		})
	}

	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&pairs).Error
}
