package gaugeMonitor

import (
	"gcore/gmonitor/internal/util"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var cache sync.Map

type gauge struct {
	gv *prometheus.GaugeVec
}

func (g *gauge) NewMonitor(labelValues ...string) GaugeMonitor {
	return newMonitor(g.gv, labelValues...)
}

// 构建方法
func NewGauge(metric string, help string, labels ...string) (g *gauge) {
	if value, ok := cache.Load(metric); ok {
		g = value.(*gauge)
	} else {
		util.Check(metric)

		g = &gauge{
			gv: prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name:      metric,
					Help:      help,
				},
				labels,
			),
		}

		prometheus.MustRegister(g.gv)
		cache.Store(metric, g)
	}

	return
}
