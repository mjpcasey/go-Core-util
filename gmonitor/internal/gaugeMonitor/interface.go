package gaugeMonitor

type GaugeVec interface {
	NewMonitor(labelValues ...string) GaugeMonitor
}

// GaugeMonitor 普通的数据指标监控
type GaugeMonitor interface {
	Add(change float64)
	Set(value float64)
	Inc() // +1
	Dec() // -1
}