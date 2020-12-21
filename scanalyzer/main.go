package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const (
	ArchiveDir = "archive"
	ReportsDir = "reports"
)

type Addr struct {
	Address	string	`xml:"addr,attr"`
	Type	string	`xml:"addrtype,attr"`
}

type Port struct {
	Id	uint16	`xml:"portid,attr"`
	Proto	string	`xml:"protocol,attr"`
}

type Host struct {
	XMLName	xml.Name	`xml:"nmaprun"`
	Address	Addr	`xml:"host>address"`
	Ports	[]Port	`xml:"host>ports>port"`
}

func main() {
	reports, err := ioutil.ReadDir(ReportsDir)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir(ArchiveDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	for _, report := range reports {
		reportPath := path.Join(ReportsDir, report.Name())
		data, err := ioutil.ReadFile(reportPath)
		if err != nil {
			log.Print(err)
			continue
		}

		host := Host{}
		err = xml.Unmarshal(data, &host)
		if err != nil {
			log.Print(err)
			continue
		}

		// TODO: Insert scan results into database
		log.Print(host)

		err = os.Rename(reportPath, path.Join(ArchiveDir, report.Name()))
		if err != nil {
			log.Print(err)
		}
	}
}
