package main

import "encoding/xml"

type Addr struct {
	Address	string	`xml:"addr,attr"`
	Type	string	`xml:"addrtype,attr"`
}

type State struct {
	State	string	`xml:"state,attr"`
}

type Port struct {
	Id	uint16	`xml:"portid,attr"`
	Proto	string	`xml:"protocol,attr"`
	State	State	`xml:"state"`
}

type Host struct {
	Address	Addr	`xml:"address"`
	Ports	[]Port	`xml:"ports>port"`
}

type Run struct {
	XMLName	xml.Name	`xml:"nmaprun"`
	Start	int64	`xml:"start,attr"`
	Hosts	[]Host	`xml:"host"`
}
