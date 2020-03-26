package codec

import (
	"gcore/gmonitor"
	"io"
	"net"
	"net/rpc"
	"time"
)

const timeListLen = 40000

// severCodec 的外层封装，添加了统计功能
type serverCodec struct {
	codec   rpc.ServerCodec
	monitor gmonitor.RequestMonitor
	timeList [timeListLen]time.Time // 存放请求的开始时间，作为后面qps rt计算
}

func (c *serverCodec) WriteResponse(resp *rpc.Response, body interface{}) (err error) {
	err = c.codec.WriteResponse(resp, body)
	if err == nil {
		c.monitor.RecordRequest("success", time.Now().Sub(c.timeList[resp.Seq%timeListLen]))
	} else {
		c.monitor.RecordRequest("failed", time.Now().Sub(c.timeList[resp.Seq%timeListLen]))
	}

	return err
}

func (c *serverCodec) ReadRequestHeader(req *rpc.Request) (err error) {

	err = c.codec.ReadRequestHeader(req)

	if err == nil {
		// 计算qps
		c.monitor.AddRequest("success")
		// 放置请求开始时间
		c.timeList[req.Seq%timeListLen] = time.Now()
	} else {
		if err != io.EOF {
			c.monitor.AddRequest("failed")
		}
	}

	return
}

func (c *serverCodec) ReadRequestBody(body interface{}) error {
	return c.codec.ReadRequestBody(body)
}

func (c *serverCodec) Close() error {
	return c.codec.Close()
}

func NewServerCodec(conn net.Conn, codecType string, monitor gmonitor.RequestMonitor) (c *serverCodec) {
	c = &serverCodec{
		codec:   NewRpcServerCodec(codecType, conn),
		monitor: monitor,
	}

	return
}
