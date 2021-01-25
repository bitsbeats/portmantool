package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/bitsbeats/portmantool/database"
)

const (
	ArchiveDir = "archive"
	ReportsDir = "reports"
)

func main() {
	db := database.InitDatabase()

	reports, err := ioutil.ReadDir(ReportsDir)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir(ArchiveDir, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

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

	time.Sleep(3 * time.Second)
}
