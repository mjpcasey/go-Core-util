package rpc

import (
	"context"
	"fmt"
	"gcore/rpc/codec"
	"github.com/xtaci/smux"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// 连接复用会话
type session struct {
	severAddr string
	conf      *ClientConfig

	session *smux.Session
	clients *clientPool

	context context.Context
	cancel  context.CancelFunc

	sync.Mutex
	sync.WaitGroup
}

// 自定义打印
func (s *session) String() string {
	if s.session != nil {
		return fmt.Sprintf("%s => %s", s.session.LocalAddr(), s.session.RemoteAddr())
	}

	return "nil"
}

// 从连接池获取 rpc.Client
func (s *session) getRpcClient() (client *rpc.Client, err error) {
	client = s.clients.get()

	if client == nil {
		client, err = s.newRpcClient()
	}

	return
}

// rpc.Client 放回连接池
func (s *session) putRpcClient(client *rpc.Client) {
	s.clients.put(client)
}

// 使用复用连接构建 rpc.Client
func (s *session) newRpcClient() (client *rpc.Client, err error) {
	s.Lock()
	stream, err := s.session.OpenStream()
	s.Unlock()

	if err == nil {
		logger.Debugf("会话 %s 新建 rpc.Client", s)

		client = rpc.NewClientWithCodec(codec.NewClientCodec(s.conf.Codec, stream))
	} else {
		logger.Errorf("会话 %s 新建 rpc.Client 失败: %s", s, err)
	}

	return
}

// 连接池维护线程
func (s *session) maintainRoutine() {
	var tempDelay = time.Second
	var maxDelay = time.Minute / 4
	for {
		var timer = time.After(tempDelay)

		select {
		case <-s.context.Done():
			logger.Debugf("会话 %s 关闭连接维护", s)
			return
		case <-timer:
			if s.session.IsClosed() {
				logger.Warnf("会话 %s 开始重建", s)

				err := s.buildSession()
				if err == nil {
					tempDelay = time.Second
				} else {
					tempDelay += time.Second * 3
					if tempDelay > maxDelay {
						tempDelay = maxDelay
					}

					logger.Errorf("会话 %s 重建失败(%s后再试): %s，", s, tempDelay, err)
				}
			} else {
				if s.clients.len() < s.conf.MuxPerConn {
					client, err := s.newRpcClient()
					if err == nil {
						s.clients.put(client)
					}
				}
			}
		}
	}
}

// 建立长连接
func (s *session) buildSession() (err error) {
	nc, err := net.DialTimeout("tcp", s.severAddr, time.Duration(s.conf.DialTimeoutSec)*time.Second)

	if err == nil {
		logger.Debugf("长连接 %s => %s 已建立", nc.LocalAddr(), nc.RemoteAddr())

		s.Lock()
		s.session, err = smux.Client(nc, nil)
		_ = s.clients.clean() // 新连接建立后把旧的客户端清理掉
		s.Unlock()
	} else {
		logger.Errorf("长连接 =>%s 建立出错: %s", s.severAddr, err)
	}

	return
}

// 启动
func (s *session) start() (err error) {
	err = s.buildSession() // 建立长连接

	if err == nil {
		go func() {
			s.Add(1)
			s.maintainRoutine() // 启动一个复用连接池维护线程
			s.Done()
		}()

		logger.Debugf("会话 %s 已启动", s)
	}

	return
}

// 关闭
func (s *session) stop() (err error) {
	s.cancel()
	s.Wait() // 等待维护线程退出

	err = s.clients.clean()
	err = s.session.Close()

	logger.Debugf("会话 %s 已关闭", s)

	return
}

// 新建连接复用会话
func newSession(serverAddr string, conf *ClientConfig) (s *session) {
	s = &session{
		severAddr: serverAddr,
		conf:      conf,
		clients:   &clientPool{},
	}
	s.context, s.cancel = context.WithCancel(context.TODO())

	return
}
