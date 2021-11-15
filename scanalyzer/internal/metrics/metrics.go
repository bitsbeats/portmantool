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
