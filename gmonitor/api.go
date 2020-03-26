package gmonitor

import (
	"database/sql"
	"gcore/gmonitor/internal/gaugeMonitor"
	"gcore/gmonitor/internal/requestMonitor"
	"gcore/gmonitor/internal/serviceMonitor/mongo"
	"gcore/gmonitor/internal/serviceMonitor/mysql"
	redisM "gcore/gmonitor/internal/serviceMonitor/redis"
	"gcore/gmonitor/internal/serviceMonitor/zookeeper"
	"github.com/globalsign/mgo"
	"github.com/go-redis/redis"
	"github.com/samuel/go-zookeeper/zk"
)

// NewRequestMonitor 创建请求监控
func NewRequestMonitor(name string) RequestMonitor {
	return requestMonitor.New(name)
}

// NewZookeeperMonitor 新建Zookeeper连接监控
func NewZookeeperMonitor(conn *zk.Conn, name string) ServiceMonitor {
	return zookeeper.New(conn, name)
}

// NewMysqlMonitor 新建Mysql监控
func NewMysqlMonitor(session *sql.DB, name string) ServiceMonitor {
	return mysql.New(session, name)
}

// NewMongoMonitor 新建Mongo监控
func NewMongoMonitor(session *mgo.Session, name string) ServiceMonitor {
	return mongo.New(session, name)
}

// NewRedisMonitor 新建Redis连接监控
func NewRedisMonitor(client *redis.Client, name string) ServiceMonitor {
	return redisM.New(client, name)
}

// NewRedisClusterMonitor 新建Redis集群连接监控
func NewRedisClusterMonitor(clusterClient *redis.ClusterClient, name string) ServiceMonitor {
	return redisM.NewCluster(clusterClient, name)
}

// NewGaugeMonitor 新建普通的数据值监控
func NewGauge(metric string, help string, labels ...string) Gauge {
	return gaugeMonitor.NewGauge(metric, help, labels...)
}
