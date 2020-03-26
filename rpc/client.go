package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gcore/gcoordinator"
	"gcore/gcoordinator/gcoordinatorTypes"
)

// 子客户端，负责一台服务器节点
type serverClient interface {
	// 服务器地址
	addr() string
	// 调用接口
	call(method string, args interface{}, reply interface{}, timeout time.Duration) (err error)
	// 启动
	start() (err error)
	// 关闭
	stop() (err error)
}

// 客户端
type client struct {
	path string // zookeeper 上服务器信息信息的读取路径
	config *ClientConfig

	severClients []serverClient // 每个 serverClient 负责一台服务器

	context context.Context
	cancel  context.CancelFunc

	indexer

	sync.RWMutex
	sync.WaitGroup
}

// Client 接口实现
func (c *client) Call(method string, args interface{}, reply interface{}) (err error) {
	for i := 1; i <= c.config.Retries; i++ {
		err = c.call(method, args, reply)

		if err == nil {
			return
		} else {
			logger.Debugf("%s 第%d/%d次调用 %s 出错: %s", c.config.Name, i, c.config.Retries, method, err)
		}
	}

	logger.Errorf("%s 调用 %s 出错: %s", c.config.Name, method, err)

	return
}
func (c *client) call(method string, args interface{}, reply interface{}) (err error) {
	serverClient, err := c.getSeverClient()

	if err == nil {
		err = serverClient.call(method, args, reply, time.Duration(c.config.CallTimeoutMs)*time.Millisecond)
	}

	return
}

// 获取一个 severeClient 用于调用
func (c *client) getSeverClient() (serverClient serverClient, err error) {
	c.RLock()
	if len(c.severClients) == 0 {
		err = ErrNoServerClient
	} else {
		serverClient = c.severClients[c.indexer.getIndex(len(c.severClients))]
	}
	c.RUnlock()

	return
}

// 监听zookeeper上的服务器更新，并调用 updateServerClients
func (c *client) updateRoutine(node gcoordinator.CoordinatorNode) {
	logger.Debugf("设置服务器节点[%s]监听", c.path)

	for {
		select {
		case <-c.context.Done():
			return
		default:
			err := node.Watch(c.context, gcoordinatorTypes.EventChildrenChanged, func(evt gcoordinatorTypes.Event) (err error) {
				logger.Debugf("监听到服务器更新: %v", evt.Data)

				err = c.updateServerClients()
				if err != nil {
					logger.Errorf("更新子客户端出错: %s", err)
				}

				return
			})

			if err != nil {
				logger.Errorf("服务器节点监听报错: %s", err)

				time.Sleep(time.Second * 10)
			}
		}
	}
}

// 更新 client.serverClients
func (c *client) updateServerClients() (err error) {
	serverAddrs, err := c.getServerAddrs()
	if err == nil {
		logger.Infof("服务器列表更新: %v", serverAddrs)
	} else {
		return
	}

	// 持有的子客户端
	var holding = make(map[string]serverClient)
	for _, serverClient := range c.severClients {
		holding[serverClient.addr()] = serverClient
	}

	// 构建子客户端或者复制现有子客户端
	var updates = make(map[string]serverClient)
	for _, serverAddr := range serverAddrs {
		updates[serverAddr] = holding[serverAddr]

		if updates[serverAddr] == nil {
			serverClient := newServerClient(serverAddr, c.config)

			err = serverClient.start()
			if err == nil {
				updates[serverAddr] = serverClient
			} else {
				logger.Errorf("构建子客户端[%s]失败: %s", serverAddr, err)

				return
			}
		}
	}

	var newClients = make([]serverClient, 0)
	for _, serverClient := range updates {
		newClients = append(newClients, serverClient)
	}

	c.Lock()
	c.severClients = newClients
	c.Unlock()

	// 关掉不复使用的子客户端
	for addr, client := range holding {
		if updates[addr] == nil {
			err = client.stop()
			if err != nil {
				logger.Warnf("移除子客户端[%s]报错: %s", addr, err)
			}
		}
	}

	return
}

// 获取服务器地址列表，优先读取配置中配置的
func (c *client) getServerAddrs() (addrs []string, err error) {
	if len(c.config.ServerAddrs) > 0 { // 优先取配置文件的服务节点
		return c.config.ServerAddrs, nil
	} else {
		return c.coordServerAddrs()
	}
}

// 读取 zookeeper 上的服务器列表
func (c *client) coordServerAddrs() (addrs []string, err error) {
	node, err := c.config.coord.Open(c.path) // 服务节点获取
	if err != nil {
		return
	}

	cNodes, err := node.GetChildrenNode()
	if err != nil {
		return nil, err
	}

	var bytes []byte
	for _, cNode := range cNodes {
		bytes, err = cNode.Get()

		if err == nil {
			var message serverInfo

			err = json.Unmarshal(bytes, &message)
			if err == nil {
				addrs = append(addrs, message.Host)
			} else {
				logger.Errorf("节点信息JSON解析出错: %s", err)

				return
			}
		} else {
			logger.Errorf("读取节点信息出错: %s", err)

			return
		}
	}

	return
}

// Client 接口实现
func (c *client) Start() (err error) {
	if len(c.config.ServerAddrs) == 0 {
		cNode, err := c.config.coord.Open(c.path)

		if err == nil {
			go func() {
				c.Add(1)
				c.updateRoutine(cNode)
				c.Done()
			}()
		} else {
			return fmt.Errorf("设置服务器列表监听失败: %s", err)
		}
	}

	err = c.updateServerClients()
	if err == nil {
		logger.Infof("%s 客户端已启动", c.config.Name)
	}

	return
}

// Client 接口实现
func (c *client) Stop() (err error) {
	c.cancel()
	c.Wait()

	c.Lock()
	for _, serverClient := range c.severClients {
		err = serverClient.stop()
	}
	c.Unlock()

	logger.Infof("%s 客户端已关闭", c.config.Name)

	return
}

func newClient(config *ClientConfig) (c *client) {
	c = &client{
		path:   fmt.Sprintf("/%s", config.Name),
		config: config,
	}
	c.context, c.cancel = context.WithCancel(context.TODO())

	return
}
func newServerClient(serverAddr string, conf *ClientConfig) serverClient {
	if conf.MuxPerConn > 0 {
		return newMuxClient(serverAddr, conf)
	}

	return newPoolClient(serverAddr, conf)
}
