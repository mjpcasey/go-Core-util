package zookeeper

import (
	"fmt"
	"gcore/gmonitor/internal/serviceMonitor/base"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
)

var (
	metric   *prometheus.GaugeVec
	initOnce sync.Once
)

// zookeeperMonitor Zookeeper连接
type monitor struct {
	base.Base
	conn *zk.Conn
}

// Ping 测试Zookeeper连接
func (m *monitor) Ping() error {
	if m.conn == nil {
		return fmt.Errorf("zookeeper 客户端 %s 不存在", m.Base.Name)
	}

	_, _, err := m.conn.Exists("/")
	return err
}

// NewZookeeperMonitor 新建Zookeeper连接监控
func New(conn *zk.Conn, name string) (m *monitor) {
	initOnce.Do(func() {
		metric = base.NewMetric("Zookeeper")
		prometheus.MustRegister(metric)
	})

	m = &monitor{
		Base: base.New(name, metric.WithLabelValues(name)),
		conn: conn,
	}

	base.StartMonitor(m)

	return
}
