package redis

import (
	"errors"
	"fmt"
	"gcore/gmonitor/internal/serviceMonitor/base"
	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var (
	metric   *prometheus.GaugeVec
	initOnce sync.Once
)

// redisMonitor Redis
type monitor struct {
	base.Base
	clusterClient *redis.ClusterClient
	client        *redis.Client
}

// Ping 测试Redis连接
func (m *monitor) Ping() error {
	if m.client == nil && m.clusterClient == nil {
		return fmt.Errorf("redis 客户端 %s 不存在", m.Base.Name)
	}

	var scmd *redis.StatusCmd
	if m.client != nil {
		scmd = m.client.Ping()
	} else {
		scmd = m.clusterClient.Ping()
	}

	if err := scmd.Err(); err != nil {
		return fmt.Errorf("redis Ping 错误 %v", err)
	} else if scmd.Val() == "PONG" {
		return nil
	} else {
		return errors.New("redis 连接未知错误")
	}
}

func New(client *redis.Client, name string) (m *monitor) {
	initOnce.Do(func() {
		metric = base.NewMetric("Redis")
		prometheus.MustRegister(metric)
	})

	m = &monitor{
		Base:   base.New(name, metric.WithLabelValues(name)),
		client: client,
	}

	base.StartMonitor(m)

	return
}

func NewCluster(client *redis.ClusterClient, name string) (m *monitor) {
	initOnce.Do(func() {
		metric = base.NewMetric("Redis")
		prometheus.MustRegister(metric)
	})

	m = &monitor{
		Base:          base.New(name, metric.WithLabelValues(name)),
		clusterClient: client,
	}

	base.StartMonitor(m)

	return
}
