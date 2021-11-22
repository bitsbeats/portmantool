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
		metrics.FailedImports.Inc()
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
		metrics.FailedImports.Inc()
		return err
	}

	err = metrics.UpdateFromDatabase(db)
	if err != nil {
		return err
	}

	return nil
}
