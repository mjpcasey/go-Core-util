package gaugeMonitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

// 普通数字监测
type monitor struct {
	gauge prometheus.Gauge
}

func (m *monitor) Add(diff float64) {
	m.gauge.Add(diff)
}

func (m *monitor) Set(value float64) {
	m.gauge.Set(value)
}

func (m *monitor) Inc() {
	m.gauge.Inc()
}

func (m *monitor) Dec() {
	m.gauge.Dec()
}

func newMonitor(gaugeVec *prometheus.GaugeVec, labels ...string) (m *monitor) {
	m = &monitor{
		gauge: gaugeVec.WithLabelValues(labels...),
	}

	return
}
