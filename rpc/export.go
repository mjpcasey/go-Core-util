package rpc

import (
	"errors"
	"fmt"
	"gcore/gcoordinator"
	"net/rpc"
)

var (
	ErrNoServerClient = errors.New("没有可用节点")
	ErrNoConnection   = errors.New("没有可用连接")
	ErrShutdown       = errors.New("连接关闭")
	ErrTimeout        = errors.New("请求超时")
	ErrEOF            = errors.New("EOF")
)

const (
	defaultMinConnPerServer = 3   // 默认单服务器最少连接数
	defaultMaxConnPerServer = 200 // 默认单服务器最大连接数，超过连接数时会每秒关闭一个连接以符合此限制
	defaultRetries          = 1   // 默认重试次数
)

// 客户端配置
type ClientConfig struct {
	Name             string   `json:"name"`             // 服务名
	MinConnPerSever  int      `json:"minConnPerSever"`  // 与每台服务器保持的最少连接数
	MaxConnPerServer int      `json:"maxConnPerServer"` // 与每台服务器保持的最少连接数
	MuxPerConn       int      `json:"muxPerConn"`       // 每条长连接的起始复用数
	Codec            string   `json:"codec"`            // 压缩解码的方式, 暂时支持gob、protobuf, 默认配置codec
	ServerAddrs      []string `json:"serverAddrs"`      // 服务节点地址
	CallTimeoutMs    int      `json:"callTimeoutMs"`    // 请求超时时间,单位:毫秒
	DialTimeoutSec   int      `json:"dialTimeoutSec"`   // 建立连接超时时间,单位:秒
	Retries          int      `json:"retries"`          // 重试次数

	coord gcoordinator.Coordinator
}

// 客户端配置初始化
func (c *ClientConfig) init(coord gcoordinator.Coordinator) *ClientConfig {
	if c.MinConnPerSever == 0 {
		c.MinConnPerSever = defaultMinConnPerServer
	}
	if c.MaxConnPerServer == 0 {
		c.MaxConnPerServer = defaultMaxConnPerServer
	}
	if c.Retries == 0 {
		c.Retries = defaultRetries
	}
	c.coord = coord

	return c
}

// 服务端配置
type ServerConfig struct {
	Name    string `json:"name"`    // 服务名称
	Addr    string `json:"addr"`    // 服务地址
	ZNode   string `json:"znode"`   // zookeeper 节点名称
	Codec   string `json:"codec"`   // 序列化类型
	ConnMux bool   `json:"connMux"` // 连接复用

	receiver interface{}
}

// 服务端配置初始化
func (c *ServerConfig) init(receiver interface{}) (*ServerConfig, error) {
	var err error

	c.Addr, err = processHost(c.Addr)
	if err != nil {
		return c, fmt.Errorf("监听配置错误: %s", err)
	}

	// zookeeper注册发现的节点名称, 没有配置的时候，直接使用本机ip+端口作为服务名称
	if c.ZNode == "" {
		c.ZNode = c.Addr
	}

	// 注册接收者
	err = rpc.Register(receiver)
	if err == nil {
		c.receiver = receiver
	}

	return c, err
}

// 客户端接口定义
type Client interface {
	Call(method string, args interface{}, reply interface{}) (err error)
	Start() (err error)
	Stop() (err error)
}

// 管理器接口定义
type Manager interface {
	NewServer(conf ServerConfig, receiver interface{}) (err error)
	NewClient(conf ClientConfig) (client Client, err error)
	Start() (err error)
	Stop() (err error)
}
