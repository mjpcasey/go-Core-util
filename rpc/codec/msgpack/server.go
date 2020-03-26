package msgpack

import (
	"bufio"
	"io"
	"net/rpc"
)

const defaultBufferSize = 4 * 1024

type serverCodec struct {
	resp ResponseHeader
	enc  *Encoder
	w    *bufio.Writer

	req RequestHeader
	dec *Decoder
	c   io.Closer
}

// NewServerCodec returns a new rpc.ServerCodec.
//
// A ServerCodec implements reading of RPC requests and writing of RPC
// responses for the server side of an RPC session. The server calls
// ReadRequestHeader and ReadRequestBody in pairs to read requests from the
// connection, and it calls WriteResponse to write a response back. The
// server calls Close when finished with the connection. ReadRequestBody
// may be called with a nil argument to force the body of the request to be
// read and discarded.
func NewServerCodec(rwc io.ReadWriteCloser) rpc.ServerCodec {
	w := bufio.NewWriterSize(rwc, defaultBufferSize)
	r := bufio.NewReaderSize(rwc, defaultBufferSize)
	return &serverCodec{
		enc: NewEncoder(w),
		w:   w,
		dec: NewDecoder(r),
		c:   rwc,
	}
}

func (c *serverCodec) WriteResponse(resp *rpc.Response, body interface{}) (err error) {
	c.resp.Method = resp.ServiceMethod
	c.resp.Seq = resp.Seq
	c.resp.Error = resp.Error

	if err = c.enc.Encode(&c.resp); err != nil {
		return err
	}
	if resp.Error == "" {
		if err = c.enc.Encode(body); err != nil {
			return err
		}
	}

	return c.w.Flush()
}

func (c *serverCodec) ReadRequestHeader(req *rpc.Request) error {
	c.req = RequestHeader{}
	if err := c.dec.Decode(&c.req); err != nil {
		return err
	}

	req.ServiceMethod = c.req.Method
	req.Seq = c.req.Seq
	return nil
}

func (c *serverCodec) ReadRequestBody(body interface{}) error {
	if body == nil {
		return nil
	}

	return c.dec.Decode(body)
}

func (c *serverCodec) Close() error { return c.c.Close() }
