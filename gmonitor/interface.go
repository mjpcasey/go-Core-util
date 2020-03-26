package gmonitor

import (
	"gcore/gmonitor/internal/gaugeMonitor"
	"time"
)

// RequestMonitor 请求监控接口
type RequestMonitor interface {
	Unregister()
	AddRequest(status string)
	RecordRequest(status string, duration time.Duration)
}

// ServiceMonitor 服务监控接口
type ServiceMonitor interface {
	Ping() error
}

type Gauge interface {
	gaugeMonitor.GaugeVec
}

// GaugeMonitor 普通的数据指标监控
type GaugeMonitor interface {
	gaugeMonitor.GaugeMonitor
}
