package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

const (
	ArchiveDir = "archive"
	ReportsDir = "reports"
)

func main() {
	db := InitDatabase()

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

		// TODO: Insert scan results into database
		log.Print(run)

		err = os.Rename(reportPath, path.Join(ArchiveDir, report.Name()))
		if err != nil {
			log.Print(err)
			continue
		}

		log.Print("done")
	}

	time.Sleep(3 * time.Second)
}
