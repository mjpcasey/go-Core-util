package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"gcore/gcoordinator"
	"gcore/gmonitor"
	"gcore/rpc/codec"
	"github.com/xtaci/smux"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// serverInfo 在zookeeper存信息，json结构
type serverInfo struct {
	Host      string // zookeeper服务节点名称，用于给client做服务链接 默认使用配置的znode配置信息，如果没有配置设置，使用ip:端口
	RealHost  string // 真实的ip和端口
	Timestamp int64  // 启动的时间戳
	HeartBeat int64  // 心跳包时间
}

// 服务器实现
type server struct {
	conf *ServerConfig

	listener net.Listener                 // tcp listener
	node     gcoordinator.CoordinatorNode // 协调器节点
	message  serverInfo                   // 在zookeeper上纪录的服务启动信息，包括运行状态，时间戳
	manager  *manager                     // rpc 管理器
	monitor  gmonitor.RequestMonitor      // 请求 计数和耗时统计

	context context.Context
	cancel  context.CancelFunc

	sync.Once
	sync.WaitGroup
}

// 连接接入监听线程
func (s *server) mainRoutine() {
	var err error

	_, port, _ := net.SplitHostPort(s.conf.Addr) // 没必要指定ip来监听
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Fatalf("%s 监听 %s 出错: %s", s.conf.Name, s.conf.Addr, err)
	}

	var tempDelay time.Duration
	for {
		conn, err := s.listener.Accept()

		if err == nil {
			logger.Infof("%s 建立会话 %s <= %s", s.conf.Name, conn.LocalAddr(), conn.RemoteAddr())

			if s.conf.ConnMux {
				go s.multiplexConn(conn)
			} else {
				go func() {
					rpc.ServeCodec(codec.NewServerCodec(conn, s.conf.Codec, s.monitor))
					logger.Infof("会话 %s 已关闭", fmt.Sprintf("%s <= %s", conn.LocalAddr(), conn.RemoteAddr()))
				}()
			}
		} else {
			select {
			case <-s.context.Done():
				logger.Debugf("%s 停止监听 %s", s.conf.Name, s.conf.Addr)
				return
			default:
				tempDelay = s.handleConnErr(tempDelay, err)
				logger.Errorf("%s 监听 %s 出错: %s", s.conf.Name, s.conf.Addr, err)
			}
		}
	}
}

// 复用长连接
func (s *server) multiplexConn(conn net.Conn) {
	var connStr = fmt.Sprintf("%s <= %s", conn.LocalAddr(), conn.RemoteAddr())

	session, err := smux.Server(conn, nil)
	if err != nil {
		logger.Warnf("%s 长连接复用 %s 出错: %s", s.conf.Name, connStr, err)
		_ = conn.Close()
		return
	}

	var sessionName = fmt.Sprintf("%s 会话 %s <=> %s", s.conf.Name, session.LocalAddr(), session.RemoteAddr())
	var streamNo = 0
	for {
		stream, err := session.AcceptStream()

		if err == nil {
			streamNo += 1
			logger.Debugf("%s 连接复用#%d", sessionName, streamNo)

			go rpc.ServeCodec(codec.NewServerCodec(stream, s.conf.Codec, s.monitor))
		} else {
			logger.Infof("%s 已关闭", sessionName)

			err = session.Close()

			return
		}
	}
}

// 处理连接接入错误
func (s *server) handleConnErr(tempDelay time.Duration, err error) time.Duration {
	// 从 net/http 抄过来的
	if ne, ok := err.(net.Error); ok && ne.Temporary() {
		if tempDelay == 0 {
			tempDelay = 5 * time.Millisecond
		} else {
			tempDelay *= 2
		}
		if max := 1 * time.Second; tempDelay > max {
			tempDelay = max
		}

		logger.Errorf("%s 临时错误: %s, %+v后重试", s.conf.Name, err, tempDelay)

		time.Sleep(tempDelay)
	}

	return tempDelay
}

// 服务器心跳线程
func (s *server) heartbeatRoutine() {
	ticker := time.Tick(time.Second)

	for {
		select {
		case <-s.context.Done():
			logger.Debugf("%s 退出心跳线程", s.conf.Name)
			return
		case <-ticker:
			err := s.node.Refresh()

			if err == nil {
				s.message.HeartBeat = time.Now().Unix()

				data, _ := json.Marshal(s.message)

				err = s.node.Set(data)
			} else {
				err = s.manager.register(s) // 尝试重新注册到zookeeper
			}

			if err != nil {
				logger.Errorf("rpc 服务[%s]心跳出错: %s", s.conf.Addr, err)
			}
		}
	}
}

// 启动
func (s *server) start() (err error) {
	s.Do(func() {
		go func() {
			s.Add(1)
			s.heartbeatRoutine()
			s.Done()
		}()

		go func() {
			s.Add(1)
			s.mainRoutine()
			s.Done()
		}()
	})

	return
}

// 关闭
func (s *server) stop() (err error) {
	s.cancel()
	err = s.listener.Close() // tcp 监听关闭
	err = s.node.Remove()
	s.Wait()

	return
}

// 新建服务器实现
func newServer(ma *manager, conf *ServerConfig) (s *server) {
	s = &server{
		conf:    conf,
		manager: ma,
		monitor: gmonitor.NewRequestMonitor(conf.Name),
	}
	s.message = serverInfo{
		Host:      s.conf.ZNode,
		RealHost:  s.conf.Addr,
		Timestamp: time.Now().Unix(),
		HeartBeat: time.Now().Unix(),
	}
	s.context, s.cancel = context.WithCancel(context.TODO())

	return
}
