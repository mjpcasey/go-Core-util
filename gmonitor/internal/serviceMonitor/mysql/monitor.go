package mysql

import (
	"database/sql"
	"fmt"
	"gcore/gmonitor/internal/serviceMonitor/base"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var (
	metric   *prometheus.GaugeVec
	initOnce sync.Once
)

// mysqlMonitor Mysql连接
type monitor struct {
	base.Base
	session *sql.DB
}

// Ping 测试Mongo连接
func (m *monitor) Ping() error {
	if m.session == nil {
		return fmt.Errorf("mysql 实例 %s 不存在", m.Name)
	}

	if err := m.session.Ping(); err != nil {
		return fmt.Errorf("mysql Ping 失败: %v", err)
	}
	return nil
}

func New(session *sql.DB, name string) (m *monitor) {
	initOnce.Do(func() {
		metric = base.NewMetric("Mysql")
		prometheus.MustRegister(metric)
	})

	m = &monitor{
		Base:    base.New(name, metric.WithLabelValues(name)),
		session: session,
	}

	base.StartMonitor(m)

	return
}
