package importer

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/bitsbeats/portmantool/scanalyzer/internal/database"
//	"github.com/bitsbeats/portmantool/scanalyzer/internal/metrics"

	"gorm.io/gorm"
)

const (
	ArchiveDir = "archive"
	ReportsDir = "reports"
)

type (
	Importer struct {
		conn *gorm.DB
	}
)

func NewImporter(db *gorm.DB) Importer {
	return Importer{db}
}

func (i Importer) importScans() {
	reports, err := ioutil.ReadDir(ReportsDir)
	if err != nil {
		log.Print(err)
		return
	}

	db := i.conn
	for _, report := range reports {
		log.Printf("Processing %s", report.Name())

		reportPath := path.Join(ReportsDir, report.Name())
		data, err := ioutil.ReadFile(reportPath)
		if err != nil {
			log.Print(err)
			continue
		}

		run := Run{}
		err = xml.Unmarshal(data, &run)
		if err != nil {
			log.Print(err)
			continue
		}

		log.Print(run)

//		rows, err := db.Model(&database.ActualState{}).Select("address, port, protocol, state").Joins("JOIN (?) ON address = a AND port = p AND protocol = proto AND scan_id = max_scan_id", db.Model(&database.ActualState{}).Select("address a, port p, protocol proto, max(scan_id) max_scan_id").Group("address, port, protocol")).Rows()
//		if err != nil {
//			log.Print(err)
//			continue
//		}

		scan := database.Scan{ID: time.Unix(run.Start, 0)}
		result := db.Create(&scan)
		if result.Error != nil {
			log.Print(result.Error)
			continue
		}

		// TODO: Needs optimization
		for _, host := range run.Hosts {
			for _, port := range host.Ports {
				state := database.ActualState{Target: database.Target{Address: host.Address.Address, Port: port.Id, Protocol: port.Proto}, State: port.State.State, ScanID: scan.ID}
				result = db.Create(&state)
				if result.Error != nil {
					log.Print(result.Error)
					continue
				}
			}
		}

		err = os.Rename(reportPath, path.Join(ArchiveDir, report.Name()))
		if err != nil {
			log.Print(err)
			continue
		}

		log.Print("done")
	}
}

func (i Importer) Run(ctx context.Context) error {
	err := os.Mkdir(ArchiveDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	go func() {
		for {
			i.importScans()

			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(3 * time.Second)
			}
		}
	}()

	return nil
}
