package zkCoordinator

import (
	"context"
	"fmt"
	"gcore/gcoordinator/gcoordinatorTypes"
	"gcore/gcoordinator/internal/interfaces"
	"gcore/gmonitor"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
	"time"
)

// zookeeperCoordinator Coordinator base on Zookeeper
type coordinator struct {
	config     *gcoordinatorTypes.Config
	timeout    time.Duration
	conn       *zk.Conn
	connEvents <-chan zk.Event

	context context.Context
	cancel  context.CancelFunc

	sync.WaitGroup
}

// init initial zookeeper module object
func (c *coordinator) init(cfg *gcoordinatorTypes.Config) (err error) {

	return nil
}

// CreateNode create a zookeeper node
func (c *coordinator) CreateNode(path string, data []byte, store bool) (interfaces.CoordinatorNode, error) {
	// 默认是临时节点
	var flags int32
	flags = zk.FlagEphemeral
	// 持久化节点
	if store {
		flags = 0
	}

	_, err := c.conn.Create(c.config.Root+path, data, flags, zk.WorldACL(zk.PermAll))
	if err != nil {
		return nil, err
	}

	node := new(node)
	if err = node.init(path, c); err != nil {
		return nil, err
	}

	return node, nil
}

func (c *coordinator) Exist(path string) (bool, error) {
	ex, _, err := c.conn.Exists(c.config.Root + path)

	if err != nil {
		return false, err
	}

	return ex, nil
}

// Open open a zookeeper node
func (c *coordinator) Open(path string) (interfaces.CoordinatorNode, error) {
	node := new(node)

	if err := node.init(path, c); err != nil {
		return nil, err
	}

	return node, nil
}

func (c *coordinator) Start() (err error) {
	c.conn, c.connEvents, err = zk.Connect(c.config.Addrs, 10*time.Second, zk.WithLogInfo(false))

	if err == nil {
		exist, _, err := c.conn.Exists(c.config.Root)

		if !exist {
			_, err = c.conn.Create(c.config.Root, nil, 0, zk.WorldACL(zk.PermAll))

			if err != nil {
				logger.Errorf("创建根节点[%s]失败: %s", c.config.Root, err)

				return err
			}
		}

		if err == nil {
			go func() {
				c.Add(1)
				c.watchConn()
				c.Done()
			}()

			gmonitor.NewZookeeperMonitor(c.conn, c.config.Root)
		}
	} else {
		return fmt.Errorf("初始化错误: %s", err)
	}

	return
}

func (c *coordinator) Stop() (err error) {
	c.cancel()
	c.Wait()

	return
}

func (c *coordinator) watchConn() {
	for {
		select {
		case <-c.connEvents:
			//logger.Debugf("连接信息: %s", evt)
		case <-c.context.Done():
			logger.Debugf("关闭连接")

			c.conn.Close()

			return
		}
	}
}

// init create and init zookeeper coordinator
func New(cfg *gcoordinatorTypes.Config) (c *coordinator) {
	c = &coordinator{
		config: cfg,
	}
	c.context, c.cancel = context.WithCancel(context.TODO())

	return
}
