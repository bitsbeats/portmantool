package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Ports = prometheus.NewCounterVec(
		prometheus.CounterOpts{
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

	FailedImports = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "portmantool_imports_failed",
		Help: "Total number of failed imports",
	})
)

func RegisterMetrics() error {
	collectors := []prometheus.Collector{Ports, RoguePorts}

	for _, c := range collectors {
		err := prometheus.Register(c)
		if err != nil {
			return err
		}
	}

	return nil
}
