package rpc

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"strconv"
	"testing"
	"time"

	tinylib_msgpack_codec "gcore/rpc/codec/tinylib_msgp"
	"gcore/rpc/test/msgp"

	. "github.com/smartystreets/goconvey/convey"
)

type MockService struct {
	name string
}

func (m *MockService) Bid(request *msgp.AdapterRequest, response *msgp.BidResponse) (err error) {
	response.ImpressionId = request.ImpId
	return nil
}

func (m *MockService) Error(request *msgp.AdapterRequest, response *msgp.BidResponse) (err error) {
	return errors.New("mock error")
}

func (m *MockService) Timeout(request *msgp.AdapterRequest, response *msgp.BidResponse) (err error) {
	time.Sleep(60 * time.Millisecond)
	response.ImpressionId = request.ImpId
	return nil
}

func startMsgpRPCServer(addr string) net.Listener {
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalf("net.Listen tcp %s, err: %v", addr, e)
	}
	rpc.Register(&MockService{})

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Print("rpc.Serve: accept:", err.Error())
				return
			}
			go func() {
				codec := tinylib_msgpack_codec.NewServerCodec(conn)
				rpc.ServeCodec(codec)
			}()
		}
	}()

	return l
}

func getMsgpClientConfig(addrs []string) *ClientConfig {
	return &ClientConfig{
		Name:            "test",
		MinConnPerSever: 10,
		Codec:           "messagepack",
		ServerAddrs:     addrs,
		CallTimeoutMs:   50,
		DialTimeoutSec:  5,
		Retries:         1,
	}
}

func TestMsgpClient(t *testing.T) {

	addr := ":7781"
	listener := startMsgpRPCServer(addr)
	defer listener.Close()
	config := getMsgpClientConfig([]string{addr})
	client := newClient(config)
	defer client.Stop()

	Convey("Msgp-RPC测试", t, func() {

		Convey("正常的RPC请求", func() {
			req := &msgp.AdapterRequest{ImpId: "1"}
			res := &msgp.BidResponse{}
			err := client.Call("MockService.Bid", req, res)
			So(err, ShouldBeNil)
			So(req.ImpId, ShouldEqual, res.ImpressionId)
		})

		Convey("不存在的服务", func() {
			req := &msgp.AdapterRequest{ImpId: "1"}
			res := &msgp.BidResponse{}
			err := client.Call("BadService.Bid", req, res)
			So(err, ShouldBeError)
			So(err.Error(), ShouldContainSubstring, "can't find service")
		})

		Convey("不存在的方法", func() {
			req := &msgp.AdapterRequest{ImpId: "1"}
			res := &msgp.BidResponse{}
			err := client.Call("MockService.BadMethod", req, res)
			So(err, ShouldBeError)
			So(err.Error(), ShouldContainSubstring, "can't find method")
		})

		Convey("错误的请求体1", func() {
			req := 111
			res := &msgp.BidResponse{}
			err := client.Call("MockService.Bid", req, res)
			So(err, ShouldBeError)
			So(err.Error(), ShouldContainSubstring, "request body invalid")
		})

		Convey("错误的请求体2", func() {
			req := &msgp.BidResponse{}
			res := &msgp.BidResponse{}
			err := client.Call("MockService.Bid", req, res)
			So(err, ShouldBeError)
			So(err.Error(), ShouldContainSubstring, "attempted to decode type")
		})

		Convey("错误的响应体1", func() {
			req := &msgp.AdapterRequest{}
			res := 222
			err := client.Call("MockService.Bid", req, res)
			So(err, ShouldBeError)
			So(err.Error(), ShouldContainSubstring, "response body invalid")
		})

		Convey("错误的响应体2", func() {
			req := &msgp.AdapterRequest{}
			res := &msgp.AdapterRequest{}
			err := client.Call("MockService.Bid", req, res)
			So(err, ShouldBeError)
			So(err.Error(), ShouldContainSubstring, "attempted to decode type")
		})

		Convey("服务端返回错误", func() {
			req := &msgp.AdapterRequest{ImpId: "1"}
			res := &msgp.BidResponse{}
			err := client.Call("MockService.Error", req, res)
			So(err, ShouldBeError)
			So(err.Error(), ShouldContainSubstring, "mock error")
		})

		Convey("请求超时(>50ms)", func() {
			req := &msgp.AdapterRequest{ImpId: "1"}
			res := &msgp.BidResponse{}
			err := client.Call("MockService.Timeout", req, res)
			So(err, ShouldBeError)
			So(err, ShouldEqual, ErrTimeout)
		})

	})
}

func BenchmarkMsgpClient(b *testing.B) {
	addr := "127.0.0.1:7780"
	listener := startMsgpRPCServer(addr)
	defer listener.Close()
	config := getMsgpClientConfig([]string{addr})
	client := newClient(config)
	defer client.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := strconv.Itoa(rand.Int())
			req := &msgp.AdapterRequest{ImpId: id}
			res := &msgp.BidResponse{}
			err := client.Call("MockService.Bid", req, res)
			if err != nil {
				b.Fatal(err)
			}
			if res.ImpressionId != req.ImpId {
				b.Fatal("数据不一致")
			}
		}
	})
}
