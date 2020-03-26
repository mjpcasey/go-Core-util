package gencode

import (
	"bufio"
	"errors"
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

func (c *serverCodec) WriteResponse(resp *rpc.Response, body interface{}) error {
	c.resp.Method = resp.ServiceMethod
	c.resp.Seq = resp.Seq
	c.resp.Error = resp.Error

	err := c.enc.Encode(&c.resp)
	if err != nil {
		return err
	}

	if resp.Error == "" {
		b, ok := body.(GenCodeType)
		if !ok {
			return errors.New("response body invalid")
		}
		if err = c.enc.Encode(b); err != nil {
			return err
		}
	}

	err = c.w.Flush()
	return err
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

	b := body.(GenCodeType)
	return c.dec.Decode(b)
}

func (c *serverCodec) Close() error { return c.c.Close() }
