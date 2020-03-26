package mongo

import (
	"fmt"
	"gcore/gmonitor/internal/serviceMonitor/base"
	"github.com/globalsign/mgo"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var (
	metric   *prometheus.GaugeVec
	initOnce sync.Once
)

// mongoMonitor Mongo连接
type monitor struct {
	base.Base
	session *mgo.Session
}

// Ping 测试Mongo连接
func (m *monitor) Ping() error {
	if m.session == nil {
		return fmt.Errorf("mongodb 实例 %s 不存在", m.Name)
	}
	if err := m.session.Ping(); err != nil {
		m.session.Refresh()
		if err := m.session.Ping(); err != nil {
			return fmt.Errorf("mongodb Ping 错误: %v", err)
		}
	}
	return nil
}

func New(session *mgo.Session, name string) (m *monitor) {
	initOnce.Do(func() {
		metric = base.NewMetric("Mongo")
		prometheus.MustRegister(metric)
	})

	m = &monitor{
		Base:    base.New(name, metric.WithLabelValues(name)),
		session: session.Copy(),
	}

	base.StartMonitor(m)

	return
}
