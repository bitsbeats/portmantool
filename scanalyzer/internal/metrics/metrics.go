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

package metrics

import (
	"github.com/bitsbeats/portmantool/scanalyzer/internal/database"
	"github.com/prometheus/client_golang/prometheus"

	"gorm.io/gorm"
)

var (
	FailedImports = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "portmantool_imports_failed_total",
		Help: "Total number of failed imports",
	})

	Ports = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "portmantool_ports",
			Help: "Number of unique host/protocol/port combinations in database (labels: host, protocol, state)",
		},
		[]string{"host", "protocol", "state"},
	)

	RoguePorts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "portmantool_ports_rogue",
			Help: "Number of ports with a state different from the expected (labels: host, protocol)",
		},
		[]string{"host", "protocol"},
	)
)

func RegisterMetrics() error {
	collectors := []prometheus.Collector{FailedImports, Ports, RoguePorts}

	for _, c := range collectors {
		err := prometheus.Register(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateFromDatabase(db *gorm.DB) error {
	state, err := database.CurrentState(db)
	if err != nil {
		return err
	}

	Ports.Reset()
	for _, s := range state {
		Ports.WithLabelValues(s.Address, s.Protocol, s.State).Inc()
	}

	diff, err := database.DiffExpected(db)
	if err != nil {
		return err
	}

	RoguePorts.Reset()
	for _, d := range diff {
		RoguePorts.WithLabelValues(d.Address, d.Protocol).Inc()
	}

	return nil
}
