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
