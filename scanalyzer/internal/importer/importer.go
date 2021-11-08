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
	"github.com/bitsbeats/portmantool/scanalyzer/internal/metrics"

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

func (i Importer) importScans() {
	reports, err := ioutil.ReadDir(ReportsDir)
	if err != nil {
		log.Print(err)
		return
	}

	for _, report := range reports {
		reportPath := path.Join(ReportsDir, report.Name())

		i.process(report.Name(), reportPath)

		_, err = os.Stat(reportPath)
		if err == nil {
			err = os.Rename(reportPath, path.Join(ArchiveDir, report.Name()))
		}
		if err != nil && !os.IsNotExist(err) {
			log.Print(err)
		}
	}

	err = metrics.UpdateFromDatabase(i.conn)
	if err != nil {
		log.Print(err)
	}
}

func (i Importer) process(report, reportPath string) {
	log.Printf("Processing %s", report)

	data, err := ioutil.ReadFile(reportPath)
	if err != nil {
		log.Print(err)
		metrics.FailedImports.Inc()
		return
	}

	run := Run{}
	err = xml.Unmarshal(data, &run)
	if err != nil {
		log.Print(err)
		metrics.FailedImports.Inc()
		return
	}

	log.Print(run)

	err = i.conn.Transaction(func(tx *gorm.DB) error {
//		rows, err := tx.Model(&database.ActualState{}).Select("address, port, protocol, state").Joins("JOIN (?) ON address = a AND port = p AND protocol = proto AND scan_id = max_scan_id", tx.Model(&database.ActualState{}).Select("address a, port p, protocol proto, max(scan_id) max_scan_id").Group("address, port, protocol")).Rows()
//		if err != nil {
//			return err
//		}

		scan := database.Scan{ID: time.Unix(run.Start, 0)}
		err := tx.Create(&scan).Error
		if err != nil {
			return err
		}

		state := make([]database.ActualState, 0, 64)
		for _, host := range run.Hosts {
			for _, port := range host.Ports {
				state = append(state, database.ActualState{
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
		err = tx.Create(&state).Error
		if err != nil {
			return err
		}

		err = os.Rename(reportPath, path.Join(ArchiveDir, report))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Print(err)
		metrics.FailedImports.Inc()
		return
	}

	log.Print("done")
}
