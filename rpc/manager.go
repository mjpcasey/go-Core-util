package rpc

import (
	"encoding/json"
	"fmt"
	"gcore/gcoordinator"
	"sync"

	"gcore/app/share"
	"gcore/glog"
)

var logger = glog.NewLogger("rpcManager")

// rpc 管理器
type manager struct {
	app     *share.AppShare
	servers map[string]map[string]*server
	coord   gcoordinator.Coordinator
	sLock   sync.RWMutex
}

// 服务端信息写入zookeeper
func (m *manager) record(server *server) error {
	path := fmt.Sprintf("/%s", server.conf.Name)

	// 创建服务根节点
	dir, err := m.coord.Open(path)
	if err != nil {
		dir, err = m.coord.CreateNode(path, []byte(path), true)
		if err != nil {
			return err
		}
	}

	// 创建子节点
	// 临时节点存在问题，可能在zookeeper上节点信息有残留，直接使用的话，后续在session超时后该临时节点会被删除
	// 所以暂时的暴力解决方法，如果发现有节点信息，先删除后新增
	// 后续在考虑优化问题
	node, err := dir.Open(server.conf.ZNode)
	if err == nil {
		err = node.Remove()
	}
	// 添加server节点，写入节点相关信息
	data, _ := json.Marshal(server.message) // 在zookeeper上纪录服务启动信息，包括运行状态，时间戳
	node, err = dir.Create(server.conf.ZNode, data)
	if err != nil {
		err = fmt.Errorf("rpc 服务注册 zookeeper 出错: %s", err)
	}

	server.node = node
	server.manager = m

	return err
}

// 更新服务端注册信息
func (m *manager) register(srv *server) (err error) {
	m.sLock.Lock()

	err = m.record(srv) // 将信息注册到zookeeper
	if err == nil {
		if m.servers[srv.conf.Name] == nil {
			m.servers[srv.conf.Name] = make(map[string]*server)
		}
		m.servers[srv.conf.Name][srv.conf.Addr] = srv

		err = srv.start() // 启动服务
	}

	m.sLock.Unlock()

	return err
}

// 构建rpc服务端
func (m *manager) NewServer(config ServerConfig, receiver interface{}) (err error) {
	conf, err := config.init(receiver)

	if err == nil {
		err = m.register(newServer(m, conf))
		if err == nil {
			logger.Infof("rpc 服务点 %s 已注册[addr=%s, znode=%s]", conf.Name, conf.Addr, conf.ZNode)
		}
	}

	return
}

// 构建rpc客户端
func (m *manager) NewClient(config ClientConfig) (client Client, err error) {
	logger.Debugf("构建客户端 %s", config.Name)

	client = newClient(config.init(m.coord))

	return
}

// 启动方法
func (m *manager) Start() (err error) {
	return
}

// 关闭方法
func (m *manager) Stop() (err error) {
	m.sLock.Lock()

	for name, servers := range m.servers {
		for addr, server := range servers {
			logger.Debugf("注销服务[name=%s, addr=%s]", name, addr)

			err = server.stop()
			if err != nil {
				logger.Errorf("注销服务[name=%s, addr=%s]出错: %s", name, addr, err)
			}

			delete(servers, addr)
		}
	}

	m.sLock.Unlock()

	return
}

// manager 构建方法
func NewManager(app *share.AppShare, co gcoordinator.Coordinator) (m *manager) {
	m = &manager{
		app:     app,
		coord:   co,
		servers: make(map[string]map[string]*server),
	}

	return m
}
