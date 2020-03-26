package base

import (
	"fmt"
	logger "gcore/glog"
	"gcore/gmonitor/internal/util"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

func NewMetric(name string) *prometheus.GaugeVec {
	util.Check(name)

	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "gmonitor",
			Name:      fmt.Sprintf("%s_up", name),
			Help:      fmt.Sprintf("The status of %s connections. Up: 1/Down: 0", name),
		},
		[]string{"name"},
	)
	return metric
}

// serviceMonitor 基础监控接口 e.q. Mysql, Mongo, Redis
type serviceMonitor interface {
	Up()
	Down()
	Ping() error
}

// StartMonitor 开始监控
func StartMonitor(monitor serviceMonitor) {
	go func() {
		for {
			if err := monitor.Ping(); err == nil {
				monitor.Up()
			} else {
				monitor.Down()
				logger.Errorf("监控失败: %v", err)
			}
			time.Sleep(time.Second * 15)
		}
	}()
}
