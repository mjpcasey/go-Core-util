package rpc

import (
	"context"
	"gcore/rpc/codec"
	"io"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// 连接池版本客户端
type poolClient struct {
	severAddr string
	config    *ClientConfig

	pool *clientPool

	context context.Context
	cancel  context.CancelFunc

	sync.Mutex
	sync.WaitGroup
}

// 接口实现
func (c *poolClient) addr() string {
	return c.severAddr
}

// 接口实现
func (c *poolClient) call(method string, args interface{}, reply interface{}, timeout time.Duration) (err error) {
	client, err := c.getRpcClient()
	if err != nil {
		return
	}

	if timeout > 0 {
		select {
		case call := <-client.Go(method, args, reply, nil).Done:
			err = call.Error
		case <-time.After(timeout):
			err = ErrTimeout
		}
	} else {
		err = client.Call(method, args, reply)
	}

	if err == nil {
		c.putRpcClient(client)
	} else {
		var abandon bool

		switch err {
		case rpc.ErrShutdown:
			err = ErrShutdown
			abandon = true
		case io.EOF:
			err = ErrEOF
			abandon = true
		default:
			if ne, ok := err.(net.Error); ok && ne.Timeout() { // 如果是超时错误，弃掉这个连接
				abandon = true
			}
		}

		if abandon {
			_ = client.Close()
		} else {
			c.putRpcClient(client)
		}
	}

	return
}

// 放回连接
func (c *poolClient) putRpcClient(client *rpc.Client) {
	c.pool.put(client)
}

// 获取一个连接
func (c *poolClient) getRpcClient() (client *rpc.Client, err error) {
	client = c.pool.get()

	if client == nil {
		client, err = c.newRpcClient()
	}

	return
}

// 新建一个连接
func (c *poolClient) newRpcClient() (client *rpc.Client, err error) {
	nc, err := net.DialTimeout("tcp", c.severAddr, time.Duration(c.config.DialTimeoutSec)*time.Second)

	if err == nil {
		logger.Infof("长连接 %s => %s 已建立", nc.LocalAddr(), nc.RemoteAddr())

		client = rpc.NewClientWithCodec(codec.NewClientCodec(c.config.Codec, nc))
	} else {
		logger.Errorf("长连接 =>%s 建立出错: %s", c.severAddr, err)
	}

	return
}

// 连接池维护线程
func (c *poolClient) maintainRoutine() {
	var tempDelay = time.Second
	var maxDelay = time.Minute / 4
	for {
		var timer = time.After(tempDelay)

		select {
		case <-c.context.Done():
			logger.Debugf("%s 关闭连接维护", c.severAddr)
			return
		case <-timer:
			if c.pool.len() < c.config.MinConnPerSever {
				client, err := c.newRpcClient()
				if err == nil {
					c.putRpcClient(client)

					tempDelay = time.Second
				} else {
					tempDelay += time.Second * 3
					if tempDelay > maxDelay {
						tempDelay = maxDelay
					}

					logger.Errorf("建立长连接 =>%s 失败(%s后再试): %s，", c.severAddr, tempDelay, err)
				}
			}

			if c.pool.len() > c.config.MaxConnPerServer {
				client := c.pool.get()
				if client != nil {
					_ = client.Close()
					logger.Infof("当前连接池大小 %s 大于偏好容量 %d，减少一个", c.pool.len(), c.config.MaxConnPerServer)
				}
			}
		}
	}
}

// 接口实现
func (c *poolClient) start() (err error) {
	client, err := c.newRpcClient() // 新建一个连接用于测试
	if err == nil {
		c.putRpcClient(client)

		go func() {
			c.Add(1)
			c.maintainRoutine() // 启动连接池维护线程
			c.Done()
		}()

		logger.Debugf("会话 %s 已启动", c)
	}

	return
}

// 接口实现
func (c *poolClient) stop() (err error) {
	c.cancel()
	c.Wait()

	err = c.pool.clean()

	logger.Debugf("会话 %s 已关闭", c)

	return
}

// 新建连接池版客户端
func newPoolClient(serverAddr string, conf *ClientConfig) (c *poolClient) {
	c = &poolClient{
		severAddr: serverAddr,
		config:    conf,
		pool:      &clientPool{},
	}
	c.context, c.cancel = context.WithCancel(context.TODO())

	return
}
