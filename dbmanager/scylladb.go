package dbmanager

import (
	"errors"
	"fmt"

	"github.com/gocql/gocql"
)

// scyllaDB logger 的实现，以打印驱动内部的 debug 信息
type gocqlLogger struct{}

func (*gocqlLogger) Print(v ...interface{}) {
	logger.Errorf("%s", fmt.Sprint(v...))
}

func (*gocqlLogger) Printf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func (*gocqlLogger) Println(v ...interface{}) {
	logger.Errorf("%s", fmt.Sprint(v...))
}

func init() {
	gocql.Logger = new(gocqlLogger)
}

// scyllaDB 封装定义
type ScyllaDBWrapper interface {
	GetNativeClient() (client *gocql.Session, err error)
	GetTableName() string
}

// scyllaDB 封装实现
type scyllaDBWrapper struct {
	Name                 string   `json:"name"`
	Hosts                []string `json:"hosts"`
	Port                 int      `json:"port"`
	Username             string   `json:"userName"`
	Password             string   `json:"password"`
	NumConns             int      `json:"numConns"`
	ConnectTimeoutSec    int      `json:"connectTimeoutSec"`
	TimeoutSec           int      `json:"timeoutSec"`
	ReconnectIntervalSec int      `json:"reconnectIntervalSec"`
	TableName            string   `json:"tableName"`

	client *gocql.Session
}

// 获取原生客户端
func (s *scyllaDBWrapper) GetNativeClient() (client *gocql.Session, err error) {
	if s.client == nil {
		err = errors.New("ScyllaDB 客户端没有初始化")
	}

	return s.client, err
}

// 获取配置表名
func (s *scyllaDBWrapper) GetTableName() string {
	return s.TableName
}
