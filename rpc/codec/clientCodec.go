package codec

import (
	"net"
	"net/rpc"
)

// clientCodec 的外层封装，添加了超时
type clientCodec struct {
	codec rpc.ClientCodec
}

func (c *clientCodec) WriteRequest(req *rpc.Request, body interface{}) (err error) {
	err = c.codec.WriteRequest(req, body)

	return
}

func (c *clientCodec) ReadResponseHeader(resp *rpc.Response) (err error) {
	err = c.codec.ReadResponseHeader(resp)

	return
}

func (c *clientCodec) ReadResponseBody(body interface{}) (err error) {
	return c.codec.ReadResponseBody(body)
}

func (c *clientCodec) Close() error {
	return c.codec.Close()
}

func NewClientCodec(codecType string, conn net.Conn) (c *clientCodec) {
	c = &clientCodec{
		codec: NewRpcClientCodec(codecType, conn),
	}

	return
}
