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

package importer

import (
	"encoding/xml"
	"log"
	"time"

	"github.com/bitsbeats/portmantool/scanalyzer/internal/database"
	"github.com/bitsbeats/portmantool/scanalyzer/internal/metrics"

	"gorm.io/gorm"
)

func Import(db *gorm.DB, data []byte) error {
	run := Run{}
	err := xml.Unmarshal(data, &run)
	if err != nil {
		fail(db)
		return err
	}

	log.Print(run)

	scan := database.Scan{
		ID:    time.Unix(run.Start, 0),
		Ports: nil,
	}
	for _, host := range run.Hosts {
		for _, port := range host.Ports {
			scan.Ports = append(scan.Ports, database.ActualState{
				Target: database.Target{
					Address: host.Address.Address,
					Port: port.Id,
					Protocol: port.Proto,
				},
				State: port.State.State,
				ScanID: scan.ID,
			})
		}
	}

	err = db.Create(&scan).Error
	if err != nil {
		fail(db)
		return err
	}

	metrics.LastImport.SetToCurrentTime()

	err = metrics.PersistToDatabase(db, database.LastImport, metrics.LastImport)
	if err != nil {
		log.Print(err)
	}

	err = metrics.UpdateFromDatabase(db)
	if err != nil {
		return err
	}

	return nil
}

func fail(db *gorm.DB) {
	metrics.FailedImports.Inc()

	err := metrics.PersistToDatabase(db, database.FailedImports, metrics.FailedImports)
	if err != nil {
		log.Print(err)
	}
}
