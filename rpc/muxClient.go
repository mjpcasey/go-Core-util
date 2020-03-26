package rpc

import (
	"io"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// 复用连接版的客户端
type muxClient struct {
	severAddr string
	config *ClientConfig
	sessions []*session
	indexer
	sync.Mutex
	sync.WaitGroup
}

func (c *muxClient) addr() string {
	return c.severAddr
}
func (c *muxClient) call(method string, args interface{}, reply interface{}, timeout time.Duration) (err error) {
	session, err := c.getSession()
	if err != nil {
		return
	}
	client, err := session.getRpcClient()
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
		session.putRpcClient(client)
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
			session.putRpcClient(client)
		}
	}

	return
}
func (c *muxClient) getSession() (session *session, err error) {
	c.Lock()
	if len(c.sessions) > 0 {
		session = c.sessions[c.indexer.getIndex(len(c.sessions))]
	} else {
		err = ErrNoConnection
	}
	c.Unlock()

	return
}
func (c *muxClient) start() (err error) {
	for i := 0; i < c.config.MinConnPerSever; i++ {
		logger.Infof("新建会话 => %s", c.severAddr)

		session := newSession(c.severAddr, c.config)

		err = session.start()
		if err == nil {
			c.sessions = append(c.sessions, session)
		} else {
			logger.Errorf("建立会话至 %s 失败: %s", c.severAddr, err)
		}
	}

	return
}
func (c *muxClient) stop() (err error) {
	for _, session := range c.sessions {
		err = session.stop()
	}

	return
}

func newMuxClient(serverAddr string, conf *ClientConfig) (c *muxClient) {
	c = &muxClient{
		severAddr: serverAddr,
		config:    conf,

		sessions: []*session{},
	}

	return
}
