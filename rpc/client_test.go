package rpc

import (
	"log"
	"net"
	"net/rpc"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var serverSleepTime time.Duration

type Handler struct {
}

type EchoReq struct {
	Message string
}

type EchoRsp struct {
	Message string
}

func (r *Handler) Echo(req *EchoReq, res *EchoRsp) error {
	if serverSleepTime > 0 {
		time.Sleep(serverSleepTime)
	}
	res.Message = req.Message
	return nil
}

type tester struct {
	name            string
	serverAddr      string
	connectionAddrs []string
	message         string
	err             error
	serverSleepTime time.Duration
	callNum         int
}

// tests 测试集
var tests = []tester{
	{name: "正常RPC请求", serverAddr: "127.0.0.1:7700", connectionAddrs: []string{"127.0.0.1:7700"}, message: "hello", err: nil, callNum: 3},
	{name: "不存在的节点请求", connectionAddrs: []string{"127.0.0.1:7702"}, message: "hello", err: ErrNoServerClient, callNum: 3},
	{name: "部分不存在的节点请求", serverAddr: "127.0.0.1:7703", connectionAddrs: []string{"127.0.0.1:7703", "127.0.0.1:8703"}, message: "hello", err: ErrNoServerClient, callNum: 3},
	{name: "请求超时(>50ms)", serverAddr: "127.0.0.1:7704", connectionAddrs: []string{"127.0.0.1:7704"}, message: "hello", err: ErrTimeout, serverSleepTime: 51 * time.Millisecond, callNum: 3},
	{name: "请求未超时(<50ms)", serverAddr: "127.0.0.1:7705", connectionAddrs: []string{"127.0.0.1:7705"}, message: "hello", err: nil, serverSleepTime: 40 * time.Millisecond, callNum: 3},
}

func TestDefaultClient(t *testing.T) {
	for _, test := range tests {
		Convey(test.name, t, func() {
			listener := startRPCServer(test.serverAddr)
			defer listener.Close()

			config := getClientConfig(test.connectionAddrs)
			client := newClient(config)
			defer client.Stop()

			var err error
			for i := 0; i < test.callNum; i++ {
				req := &EchoReq{
					Message: test.message,
				}
				res := &EchoRsp{}
				serverSleepTime = test.serverSleepTime
				err = client.Call("Handler.Echo", req, res)
				if err != nil {
					So(err, ShouldEqual, test.err)
					continue
				}

				So(req.Message, ShouldEqual, res.Message)
			}
		})
	}
}

func BenchmarkDefaultClient(b *testing.B) {
	addr := ":4399"
	listener := startRPCServer(addr)
	defer listener.Close()

	config := getClientConfig([]string{addr})
	client:= newClient(config)
	defer client.Stop()

	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := &EchoReq{
			Message: "ok",
		}
		res := &EchoRsp{}
		err = client.Call("Handler.Echo", req, res)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func getClientConfig(addrs []string) *ClientConfig {
	return &ClientConfig{
		Name:            "test",
		MinConnPerSever: 10,
		Codec:           "",
		ServerAddrs:     addrs,
		CallTimeoutMs:   50,
		DialTimeoutSec:  5,
		Retries:         1,
	}
}

func startRPCServer(addr string) net.Listener {
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalf("net.Listen tcp %s, err: %v", addr, e)
	}
	rpc.Register(new(Handler))
	go rpc.Accept(l)
	return l
}
